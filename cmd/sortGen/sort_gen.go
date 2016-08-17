package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	GoFile    = "GOFILE"
	GoPackage = "GOPACKAGE"
)

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

	matches, err := getMatches(goFile)
	if err != nil {
		panic(err)
	}

	err = genMatches(matches, goPkg, wd)
	if err != nil {
		panic(err)
	}
}

// TODO add support for
//
// 1. support slices of imported types (would mean match.typ could be different)
// 2. support for orderers with errors?

func genMatches(matches map[string][]string, pkg string, path string) error {
	for typ, funs := range matches {
		ofName := filepath.Join(path, "gen_"+typ+"_sorter.go")

		out := bytes.NewBuffer([]byte(`package ` + pkg + `

		import "sort"
		import "github.com/myitcv/sorter"

		`))

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

	if strings.HasPrefix(orderFn, "O") {
		split = "Order"
	} else {
		lower = true
		split = "order"
	}

	parts := strings.SplitAfterN(orderFn, split, 2)

	if lower {
		return "sort" + parts[1]
	} else {
		return "Sort" + parts[1]
	}
}

func getMatches(file string) (map[string][]string, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	theImport := ""

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

	if theImport == "" {
		return nil, nil
	}

	matches := make(map[string][]string)

Decls:
	for _, d := range f.Decls {
		fun, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fn := fun.Name.Name

		lower := false
		upper := false

		// TODO check that the fn is not just order or Order
		if strings.HasPrefix(fn, "order") {
			lower = true
		} else if strings.HasPrefix(fn, "Order") {
			upper = true
		}

		if !lower && !upper {
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

		if typ.Sel.Name != "Order" {
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

		// TODO we can be looser here... can pretty much allow
		// anything because we know we have a slice
		sliceIdent, ok := at.Elt.(*ast.Ident)
		if !ok {
			continue
		}

		for i := 1; i < len(paramList); i++ {
			if id, ok := paramList[i].(*ast.Ident); !ok || id.Name != "int" {
				continue Decls
			}
		}

		funs, ok := matches[sliceIdent.Name]
		if !ok {
			funs = make([]string, 0)
		}
		funs = append(funs, fun.Name.Name)
		matches[sliceIdent.Name] = funs
	}

	return matches, nil
}
