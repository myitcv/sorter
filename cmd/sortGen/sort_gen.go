// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/imports"

	"github.com/myitcv/gogenerate"
	"github.com/myitcv/sorter"
)

const (
	sortGenCmd  = "sortGen"
	orderPrefix = "order"
)

// matching related vars
var (
	orderFnRegex    *regexp.Regexp
	lowerOrder      string
	upperOrder      string
	invalidFileChar *regexp.Regexp
)

// flags
var (
	fLicenseFile = gogenerate.LicenseFileFlag()
	fGoGenLog    = gogenerate.LogFlag()
)

func init() {
	r, n := utf8.DecodeRuneInString(orderPrefix)
	if r == utf8.RuneError {
		fatalf("OrderPrefix not a UTF8 string?")
	}

	l := string(unicode.ToLower(r))
	u := string(unicode.ToUpper(r))

	suffix := orderPrefix[n:]

	lowerOrder = l + suffix
	upperOrder = u + suffix

	orderFunctionPattern := `^[` + l + u + `]` + suffix + `[[:word:]]+`
	orderFnRegex = regexp.MustCompile(orderFunctionPattern)

	invalidFileChar = regexp.MustCompile(`[[:^word:]]`)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(sortGenCmd + ": ")

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	flag.Parse()

	gogenerate.DefaultLogLevel(fGoGenLog, gogenerate.LogFatal)

	envFileName, ok := os.LookupEnv(gogenerate.GOFILE)
	if !ok {
		fatalf("env not correct; missing %v", gogenerate.GOFILE)
	}

	envPkgName, ok := os.LookupEnv(gogenerate.GOPACKAGE)
	if !ok {
		fatalf("env not correct; missing %v", gogenerate.GOPACKAGE)
	}

	wd, err := os.Getwd()
	if err != nil {
		fatalf("unable to get working directory: %v", err)
	}

	// are we running against the first file that contains the sortGen directive?
	// if not return
	dirFiles, err := gogenerate.FilesContainingCmd(wd, sortGenCmd)
	if err != nil {
		fatalf("could not determine if we are the first file: %v", err)
	}

	if len(dirFiles) == 0 {
		fatalf("cannot find any files containing the %v directive", sortGenCmd)
	}

	if envFileName != dirFiles[0] {
		return
	}

	// if we get here, we know we are the first file...

	matches := getMatchesForPkg(wd, envPkgName)

	licenseHeader, err := gogenerate.CommentLicenseHeader(fLicenseFile)
	if err != nil {
		fatalf("could not comment license file: %v", err)
	}

	genMatches(matches, envPkgName, wd, licenseHeader)
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

// getMatchesForPkg returns a map[string]fileMatches where the string is the file path where
// the matches were found
func getMatchesForPkg(path string, pkgName string) map[string]fileMatches {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		fatalf("could not parse dir %v: %v", path, err)
	}

	pkg, ok := pkgs[pkgName]
	if !ok {
		// TODO come up with a proper error strategy
		panic("Oh dear...")
	}

	files := make(map[*ast.File]string)

	for fn, f := range pkg.Files {
		// we can safely skip files that this generator generated
		if gogenerate.FileGeneratedBy(fn, sortGenCmd) {
			continue
		}

		// now see whether it imports the sorter package
		sorterImport := ""

		for _, s := range f.Imports {
			if s.Path.Value == `"`+sorter.PkgName+`"` {
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
		return nil
	}

	realRes := make(map[string]fileMatches)

	for f, theImport := range files {
		matches := getMatchesFromFile(f, fset, theImport, pkgName)

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

		if len(funs) > 0 {
			realRes[fileName] = fileMatches{
				funs:    funs,
				imports: importMap,
			}
		}
	}

	return realRes
}

type match struct {
	// the actual function/method that has matched
	fun *ast.FuncDecl

	// the "type" of the slice parameter (the first one); i.e.
	// the expression that appears after the '[]'
	orderTyp ast.Expr
}

func getMatchesFromFile(f *ast.File, fset *token.FileSet, theImport string, goPkg string) []match {
	var matches []match

Decls:
	for _, d := range f.Decls {
		fun, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fn := fun.Name.Name

		if !orderFnRegex.MatchString(fn) {
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

		if typ.Sel.Name != sorter.OrderedName {
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

		infof("found a match at %v", fset.Position(fun.Pos()))

		matches = append(matches, match{
			fun:      fun,
			orderTyp: at.Elt,
		})
	}

	return matches
}

func genMatches(matches map[string]fileMatches, pkg string, path string, licenseHeader string) {

	license := licenseHeader

	// we need to generate one file for non-test matches... and one for test matches

	for fn, fm := range matches {
		buf := bytes.NewBuffer(nil)

		name := filepath.Base(fn)

		if strings.HasSuffix(name, "_test.go") {
			name = strings.TrimSuffix(name, "_test.go")
		} else {
			name = strings.TrimSuffix(name, ".go")
		}

		name = gogenerate.NameFile(name, sortGenCmd)
		ofName := filepath.Join(path, name)

		buf.WriteString(license)

		buf.WriteString(`// File generated by sortGen - do not edit

		`)

		buf.WriteString(`package ` + pkg + `

			import "sort"
			import "` + sorter.PkgName + `"

		`)

		for i := range fm.imports {
			fmt.Fprintln(buf, "import", i)
		}

		for _, toGen := range fm.funs {
			sortName, stableName := sortFunctions(toGen.name)

			x := ""

			if toGen.recv != "" {
				x = toGen.recvVar + "."
			}

			_, err := fmt.Fprint(buf, `
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

			if err != nil {
				fatalf("unable to print template out: %v", err)
			}
		}

		toWrite := buf.Bytes()

		res, err := imports.Process(ofName, toWrite, nil)
		if err == nil {
			toWrite = res
		}

		wrote, err := gogenerate.WriteIfDiff(toWrite, ofName)
		if err != nil {
			fatalf("could not write %v: %v", ofName, err)
		}

		if wrote {
			infof("writing %v", ofName)
		} else {
			infof("skipping writing of %v; it's identical", ofName)
		}
	}
}

func sortFunctions(orderFn string) (string, string) {
	// TODO this can be improved

	lower := false
	split := ""

	if strings.HasPrefix(orderFn, upperOrder) {
		split = upperOrder
	} else {
		lower = true
		split = lowerOrder
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

func fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func infof(format string, args ...interface{}) {
	if *fGoGenLog == string(gogenerate.LogInfo) {
		log.Printf(format, args...)
	}
}
