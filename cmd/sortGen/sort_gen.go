// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
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

	"github.com/myitcv/sorter/gen"
)

const (
	_SorterPackage = "github.com/myitcv/sorter"

	_OrderPrefix = "order"

	_SorterOrderTypeName = "Ordered"

	goFileEnv = "GOFILE"
	goPkgEnv  = "GOPACKAGE"

	// TODO shouldn't hard-code this to sortGen, should use os.Arg(0)?
	_GoGenPattern = `^//go:generate +sortGen`
)

var _GoGenerateRegex *regexp.Regexp
var _OrderFunctionRegex *regexp.Regexp
var _LowerOrder string
var _UpperOrder string
var _InvalidFileChar *regexp.Regexp

var (
	fLicenseFile = flag.String("licenseFile", "", "file containing an uncommented license header")
	fGoGenLog    = flag.String("gglog", "fatal", "log level; one of info, warning, error, fatal")
)

const (
	LogInfo    = "info"
	LogWarning = "warning"
	LogError   = "error"
	LogFatal   = "fatal"
)

var errNotFirstFile = errors.New("Not first go generate file")

func init() {
	r, n := utf8.DecodeRuneInString(_OrderPrefix)
	if r == utf8.RuneError {
		panic("OrderPrefix not a UTF8 string?")
	}

	l := string(unicode.ToLower(r))
	u := string(unicode.ToUpper(r))

	suffix := _OrderPrefix[n:]

	_LowerOrder = l + suffix
	_UpperOrder = u + suffix

	orderFunctionPattern := `^[` + l + u + `]` + suffix + `[[:word:]]+`
	_OrderFunctionRegex = regexp.MustCompile(orderFunctionPattern)

	_GoGenerateRegex = regexp.MustCompile(_GoGenPattern)

	_InvalidFileChar = regexp.MustCompile(`[[:^word:]]`)
}

func main() {
	flag.Parse()

	if *fGoGenLog == "" {
		*fGoGenLog = LogFatal
	}

	fmt.Printf("sortGen os.Args: %v, %#v", len(os.Args), os.Args)

	envFile, ok := os.LookupEnv(goFileEnv)
	if !ok {
		panic("Env not correct; missing " + goFileEnv)
	}

	envPkg, ok := os.LookupEnv(goPkgEnv)
	if !ok {
		panic("Env not correct; missing " + goPkgEnv)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	matches, err := getMatchesForPkg(wd, envFile, envPkg)
	if err != nil {
		if err == errNotFirstFile {
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

	licenseHeader := ""

	if *fLicenseFile != "" {
		byts, err := ioutil.ReadFile(*fLicenseFile)
		if err != nil {
			panic(err)
		}

		licenseHeader = string(byts)
	}

	err = genMatches(matches, envPkg, wd, licenseHeader)
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
		if !fileNotGenerated(e) {
			err = os.Remove(filepath.Join(dir, e.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fileNotGenerated(file os.FileInfo) bool {
	fn := file.Name()
	return !strings.HasPrefix(fn, gen.GenFilePrefix) || !strings.HasSuffix(fn, gen.GenFileSuffix)
}

type fileMatches struct {
	// the string is the import name and quoted path combined
	imports map[string]bool

	funs []sortFunToGen
}

type sortFunToGen struct {
	name    string
	recv    string
	recvVar string
	typ     string
}

// getMatchesForPkg returns a map[string]fileMatches where the string is the filename where
// the matches were found
func getMatchesForPkg(path string, envFile string, envPkg string) (map[string]fileMatches, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, fileNotGenerated, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	pkg, ok := pkgs[envPkg]
	if !ok {
		// TODO come up with a proper error strategy
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
				if _GoGenerateRegex.MatchString(com.Text) {
					foundComment = true
					break FileComments
				}
			}
		}

		if !foundComment {
			continue
		}

		// now see whether it imports the sorter package
		sorterImport := ""

		for _, s := range f.Imports {
			if s.Path.Value == `"`+_SorterPackage+`"` {
				if s.Name != nil {
					sorterImport = s.Name.Name
				} else {
					naked := strings.Trim(s.Path.Value, `"`)
					parts := strings.Split(naked, "/")
					sorterImport = parts[len(parts)-1]
				}
			}
		}

		if sorterImport != "" {
			files[f] = sorterImport
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

		if fileList[0] != envFile {
			return nil, errNotFirstFile
		}
	}

	realRes := make(map[string]fileMatches)

	for f, theImport := range files {
		matches, err := getMatchesFromFile(f, fset, theImport, envPkg)
		if err != nil {
			return nil, err
		}

		var funs []sortFunToGen
		importMap := make(map[string]bool)

		// we need to union the list of functions
		for _, match := range matches {
			var buf bytes.Buffer

			err := printer.Fprint(&buf, fset, match.orderTyp)
			if err != nil {
				// TODO
				panic(err)
			}

			sliceIdent := buf.String()

			recv := ""
			recvVar := ""

			if match.fun.Recv != nil {
				var buf bytes.Buffer

				// we know at this point we have a valid method...
				recvVar = match.fun.Recv.List[0].Names[0].Name

				buf.WriteString("(")
				buf.WriteString(recvVar)

				// TODO should handle error here because it's not stated the
				// print only supports type X, it's implementation detail
				err := printer.Fprint(&buf, fset, match.fun.Recv.List[0].Type)
				if err != nil {
					panic(err)
				}

				buf.WriteString(")")

				recv = buf.String()
			}

			// we need to calculate the required imports
			importMatches := findImports(match.orderTyp, f.Imports)

			for i := range importMatches {
				importName := i.Path.Value
				if i.Name != nil {
					importName = i.Name.Name + " " + importName
				}

				importMap[importName] = true
			}

			funs = append(funs, sortFunToGen{
				name:    match.fun.Name.Name,
				typ:     sliceIdent,
				recv:    recv,
				recvVar: recvVar,
			})
		}

		fileName := fset.Position(f.Pos()).Filename
		basename := strings.TrimSuffix(filepath.Base(fileName), ".go")

		realRes[basename] = fileMatches{
			funs:    funs,
			imports: importMap,
		}
	}

	return realRes, nil
}

type match struct {
	// the actual function/method that has matched
	fun *ast.FuncDecl

	// the "type" of the slice parameter (the first one); i.e.
	// the expression that appears after the '[]'
	orderTyp ast.Expr
}

func getMatchesFromFile(f *ast.File, fset *token.FileSet, theImport string, goPkg string) ([]match, error) {
	var matches []match

Decls:
	for _, d := range f.Decls {
		fun, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fn := fun.Name.Name

		if !_OrderFunctionRegex.MatchString(fn) {
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

		if typ.Sel.Name != _SorterOrderTypeName {
			continue
		}

		if fun.Type.Params == nil {
			continue
		}

		// we need to gather the number of params....
		var paramList []ast.Expr
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

		matches = append(matches, match{
			fun:      fun,
			orderTyp: at.Elt,
		})
	}

	return matches, nil
}

// TODO add support for
//
// 1. support for orderers with errors?

func genMatches(matches map[string]fileMatches, pkg string, path string, licenseHeader string) error {

	license := commentString(licenseHeader)

	for file, fm := range matches {
		name := "gen_" + file + ".sortGen.go"
		ofName := filepath.Join(path, name)

		out := bytes.NewBuffer(nil)

		out.WriteString(license)

		out.WriteString(`package ` + pkg + `

		import "sort"
		import "` + _SorterPackage + `"

		`)

		for i := range fm.imports {
			fmt.Fprintln(out, "import", i)
		}

		for _, toGen := range fm.funs {
			sortName, stableName := sortFunctions(toGen.name)

			x := ""

			if toGen.recv != "" {
				x = toGen.recvVar + "."
			}

			fmt.Fprint(out, `
			func `+toGen.recv+` `+sortName+`(vs []`+toGen.typ+`) {
				sort.Sort(&sorter.Wrapper{
					LenFunc: func() int {
						return len(vs)
					},
					LessFunc: func(i, j int) bool {
						return bool(`+x+toGen.name+`(vs, i, j))
					},
					SwapFunc: func(i, j int) {
						vs[i], vs[j] = vs[j], vs[i]
					},
				})
			}
			func `+toGen.recv+` `+stableName+`(vs []`+toGen.typ+`) {
				sort.Sort(&sorter.Wrapper{
					LenFunc: func() int {
						return len(vs)
					},
					LessFunc: func(i, j int) bool {
						return bool(`+x+toGen.name+`(vs, i, j))
					},
					SwapFunc: func(i, j int) {
						vs[i], vs[j] = vs[j], vs[i]
					},
				})
			}
			`)
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

func sortFunctions(orderFn string) (string, string) {
	// TODO this can be improved

	lower := false
	split := ""

	if strings.HasPrefix(orderFn, _UpperOrder) {
		split = _UpperOrder
	} else {
		lower = true
		split = _LowerOrder
	}

	parts := strings.SplitAfterN(orderFn, split, 2)

	if lower {
		return "sort" + parts[1], "stableSort" + parts[1]
	}

	return "Sort" + parts[1], "StableSort" + parts[1]
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

func commentLicense(licenseFile string) (string, error) {
	if licenseFile == "" {
		return "", nil
	}

	file, err := os.Open(licenseFile)
	if err != nil {
		return "", err
	}

	res := ""

	lastLineEmpty := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			lastLineEmpty = true
		}
		res = res + fmt.Sprintln("//", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("License file %v, %v\n", licenseFile, err)
		return "", err
	}

	// ensure we have a space before package
	if !lastLineEmpty {
		res = res + "\n"
	}

	return res, nil
}

func commentString(r string) string {
	res := ""

	buf := bytes.NewBuffer([]byte(r))

	lastLineEmpty := false
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			lastLineEmpty = true
		}
		res = res + fmt.Sprintln("//", line)
	}

	if err := scanner.Err(); err != nil {
		// this really would be exceptional... because we passed in a string
		panic(err)
	}

	// ensure we have a space before package
	if !lastLineEmpty {
		res = res + "\n"
	}

	return res
}
