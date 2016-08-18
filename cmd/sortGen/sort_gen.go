package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	OrderPrefix = "order"

	SorterOrderTypeName = "Order"

	GoFile        = "GOFILE"
	GoPackage     = "GOPACKAGE"
	GenFilePrefix = "gen_"
	GenFileSuffix = "_sorter.go"

	// TODO shouldn't hard-code this to sortGen, should use os.Arg(0)?
	GoGenPattern = `^//go:generate +sortGen`
)

var GoGenerateRegex *regexp.Regexp
var OrderFunctionRegex *regexp.Regexp
var LowerOrder string
var UpperOrder string
var InvalidFileChar *regexp.Regexp

func init() {
	r, n := utf8.DecodeRuneInString(OrderPrefix)
	if r == utf8.RuneError {
		panic("OrderPrefix not a UTF8 string?")
	}

	l := string(unicode.ToLower(r))
	u := string(unicode.ToUpper(r))

	suffix := OrderPrefix[n:]

	LowerOrder = l + suffix
	UpperOrder = u + suffix

	orderFunctionPattern := `^[` + l + u + `]` + suffix + `[[:word:]]+`
	OrderFunctionRegex = regexp.MustCompile(orderFunctionPattern)

	GoGenerateRegex = regexp.MustCompile(GoGenPattern)

	InvalidFileChar = regexp.MustCompile(`[[:^word:]]`)
}

var NotFirstFile = errors.New("Not first go generate file")

func main() {
	goFile, ok := os.LookupEnv(GoFile)
	if !ok {
		panic("Env not correct; missing " + GoFile)
	}

	goPkg, ok := os.LookupEnv(GoPackage)
	if !ok {
		panic("Env not correct; missing " + GoPackage)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	matches, err := getMatchesForPkg(wd, goFile, goPkg)
	if err != nil {
		if err == NotFirstFile {
			return
		}

		panic(err)
	}

	// if we get here, we know we are the first file... hence
	// we can safely delete existing generated files before
	// generating new ones
	err = removeGeneratedFiles(wd)
	if err != nil {
		panic(err)
	}

	err = genMatches(matches, goPkg, wd)
	if err != nil {
		panic(err)
	}
}

func removeGeneratedFiles(dir string) error {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		fn := e.Name()
		if strings.HasPrefix(fn, GenFilePrefix) && strings.HasSuffix(fn, GenFileSuffix) {
			err = os.Remove(filepath.Join(dir, fn))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func filterGeneratedFiles(file os.FileInfo) bool {
	fn := file.Name()
	return !strings.HasPrefix(fn, GenFilePrefix) || !strings.HasSuffix(fn, GenFileSuffix)
}

type fileMatches struct {
	imports map[string]bool
	funs    map[string][]string
}

func getMatchesForPkg(dir string, goFile string, goPkg string) (map[string]fileMatches, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, filterGeneratedFiles, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	theImport := ""

	pkg, ok := pkgs[goPkg]
	if !ok {
		panic("Oh dear...")
	}

	files := make(map[*ast.File]string)

	for _, f := range pkg.Files {
		cm := ast.NewCommentMap(fset, f, f.Comments)

		foundComment := false

		// if we find a comment that's great
	FileComments:
		for _, cg := range cm[f] {
			for _, com := range cg.List {
				if GoGenerateRegex.MatchString(com.Text) {
					foundComment = true
					break FileComments
				}
			}
		}

		if !foundComment {
			continue
		}

		// now see whether it imports the sorter package
		theImport = ""

		for _, s := range f.Imports {
			if s.Path.Value == `"github.com/myitcv/sorter"` {
				if s.Name != nil {
					theImport = s.Name.Name
				} else {
					naked := strings.Trim(s.Path.Value, `"`)
					parts := strings.Split(naked, "/")
					theImport = parts[len(parts)-1]
				}
			}
		}

		if theImport != "" {
			files[f] = theImport
		}
	}

	if len(files) == 0 {
		return nil, nil
	} else if len(files) > 1 {
		// we need to ascertain whether the file we have been called for
		// is the first in the list - this logic depends on the defined
		// behaviour of go generate (see go generate --help)
		var fileList []string
		for f := range files {
			fn := filepath.Base(fset.Position(f.Pos()).Filename)
			fileList = append(fileList, fn)
		}

		sort.Sort(sort.StringSlice(fileList))

		if fileList[0] != goFile {
			return nil, NotFirstFile
		}
	}

	realRes := make(map[string]fileMatches)

	for f, theImport := range files {
		res, err := getMatchesFromFile(f, fset, theImport, goPkg)
		if err != nil {
			return nil, err
		}

		typeMap := make(map[string][]string)
		importMap := make(map[string]bool)

		// we need to union the list of functions
		for typ, fns := range res {
			// TODO we need to do more here... because we need to work out
			// whether there should be additional imports in the generated
			// files
			var buf bytes.Buffer
			printer.Fprint(&buf, fset, typ)
			sliceIdent := buf.String()

			// we need to calculate the required imports
			importMatches := findImports(typ, f.Imports)

			for i := range importMatches {
				importName := i.Path.Value
				if i.Name != nil {
					importName = i.Name.Name + " " + importName
				}

				importMap[importName] = true
			}

			if typFns, ok := typeMap[sliceIdent]; ok {
				typFns = append(typFns, fns...)
				typeMap[sliceIdent] = typFns
			} else {
				typeMap[sliceIdent] = fns
			}
		}

		fileName := fset.Position(f.Pos()).Filename
		basename := strings.TrimSuffix(filepath.Base(fileName), ".go")

		matches := fileMatches{
			funs:    typeMap,
			imports: importMap,
		}

		realRes[basename] = matches
	}

	return realRes, nil
}

func getMatchesFromFile(f *ast.File, fset *token.FileSet, theImport string, goPkg string) (map[ast.Expr][]string, error) {
	matches := make(map[ast.Expr][]string)

Decls:
	for _, d := range f.Decls {
		fun, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fn := fun.Name.Name

		if !OrderFunctionRegex.MatchString(fn) {
			continue
		}

		if fun.Type.Results == nil || len(fun.Type.Results.List) != 1 {
			continue
		}

		typ, ok := fun.Type.Results.List[0].Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		id, ok := typ.X.(*ast.Ident)
		if !ok {
			continue
		}

		if id.Name != theImport {
			continue
		}

		if typ.Sel.Name != SorterOrderTypeName {
			continue
		}

		if fun.Type.Params == nil {
			continue
		}

		// we need to gather the number of params....
		paramList := make([]ast.Expr, 0)
		for _, f := range fun.Type.Params.List {
			for _ = range f.Names {
				paramList = append(paramList, f.Type)
			}
		}

		if len(paramList) != 3 {
			continue
		}

		at, ok := paramList[0].(*ast.ArrayType)
		if !ok || at.Len != nil {
			continue
		}

		for i := 1; i < len(paramList); i++ {
			if id, ok := paramList[i].(*ast.Ident); !ok || id.Name != "int" {
				continue Decls
			}
		}

		funs, ok := matches[at.Elt]
		if !ok {
			funs = make([]string, 0)
		}
		funs = append(funs, fun.Name.Name)
		matches[at.Elt] = funs
	}

	return matches, nil
}

// TODO add support for
//
// 1. support slices of imported types (would mean match.typ could be different)
// 2. support for orderers with errors?

func genMatches(matches map[string]fileMatches, pkg string, path string) error {
	for typ, funs := range matches {
		name := "gen_" + typ + "_sorter.go"
		ofName := filepath.Join(path, name)

		out := bytes.NewBuffer([]byte(`package ` + pkg + `

		import "sort"
		import "github.com/myitcv/sorter"

		`))

		for i := range funs.imports {
			fmt.Fprintln(out, "import", i)
		}

		for typ, funs := range funs.funs {
			for _, fun := range funs {
				sortName := sortFunction(fun)

				fmt.Fprint(out, `
				func `+sortName+`(vs []`+typ+`) {
					sort.Sort(&sorter.Wrapper{
						LenFunc: func() int {
							return len(vs)
						},
						LessFunc: func(i, j int) bool {
							return bool(`+fun+`(vs, i, j))
						},
						SwapFunc: func(i, j int) {
							vs[i], vs[j] = vs[j], vs[i]
						},
					})
				}
				`)
			}
		}

		toWrite := out.Bytes()

		res, err := format.Source(toWrite)
		if err == nil {
			toWrite = res
		}

		of, err := os.Create(ofName)
		if err != nil {
			return err
		}

		_, err = of.Write(toWrite)
		of.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func sortFunction(orderFn string) string {
	// TODO this can be improved

	lower := false
	split := ""

	if strings.HasPrefix(orderFn, UpperOrder) {
		split = UpperOrder
	} else {
		lower = true
		split = LowerOrder
	}

	parts := strings.SplitAfterN(orderFn, split, 2)

	if lower {
		return "sort" + parts[1]
	} else {
		return "Sort" + parts[1]
	}
}

func typeStringToFileName(s string) string {
	// safely translate to a filename
	res := s
	res = strings.Replace(res, "[]", "sl_", -1)
	res = strings.Replace(res, "*", "p_", -1)

	return InvalidFileChar.ReplaceAllString(res, "_")
}

type importFinder struct {
	imports []*ast.ImportSpec
	matches map[*ast.ImportSpec]bool
}

func (i *importFinder) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.SelectorExpr:
		if x, ok := node.X.(*ast.Ident); ok {
			for _, imp := range i.imports {
				if imp.Name != nil {
					if x.Name == imp.Name.Name {
						i.matches[imp] = true
					}
				} else {
					cleanPath := strings.Trim(imp.Path.Value, "\"")
					parts := strings.Split(cleanPath, "/")
					if x.Name == parts[len(parts)-1] {
						i.matches[imp] = true
					}
				}
			}

		}
	}

	return i
}

func findImports(exp ast.Expr, imports []*ast.ImportSpec) map[*ast.ImportSpec]bool {
	finder := &importFinder{
		imports: imports,
		matches: make(map[*ast.ImportSpec]bool),
	}

	ast.Walk(finder, exp)

	return finder.matches
}
