// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {
	matches, err := getMatchesForPkg("_testFiles/", "to_parse.go", "main")

	if err != nil {
		t.Fatalf("Expected err to be nil, got %v", err)
	}

	// number of types matched
	if len(matches) != 2 {
		t.Fatalf("We got %v matches instead of 2", len(matches))
	}

	funs := matches["to_parse"]

	if len(funs.funs) != 5 {
		t.Fatalf("We got %v function matches instead of 5", len(funs.funs))
	}

	tmpDir, err := ioutil.TempDir("", "sortGen_temp")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	err = genMatches(matches, "main", tmpDir)

	if err != nil {
		t.Fatalf("Expected gen err to be nil, got %v", err)
	}

	checkFiles, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	for _, f := range checkFiles {
		fmt.Printf(">> %v\n", f.Name())
		outFile, err := ioutil.ReadFile(filepath.Join(tmpDir, f.Name()))
		if err != nil {
			panic(err)
		}

		fmt.Println(string(outFile))
	}
}
