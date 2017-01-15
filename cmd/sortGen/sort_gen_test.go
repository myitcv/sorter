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
	fileMatches := getMatchesForPkg("_testFiles/", "to_parse.go", "main")

	for k, v := range fileMatches {
		fmt.Printf("%v: %v\n", k, v)
	}

	// number of types matched
	if len(fileMatches) != 2 {
		t.Fatalf("We got %v matches instead of 2", len(fileMatches))
	}

	funs := fileMatches["_testFiles/to_parse.go"]

	if len(funs.funs) != 5 {
		t.Fatalf("We got %v function matches instead of 5", len(funs.funs))
	}

	tmpDir, err := ioutil.TempDir("", "sortGen_temp")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpLicenseFileFi, err := ioutil.TempFile(os.TempDir(), "sortGet_tempLicenseFile")
	if err != nil {
		panic(err)
	}

	_, err = tmpLicenseFileFi.WriteString("My favourite license")
	if err != nil {
		panic(err)
	}

	tmpLicenseFileFi.Close()

	genMatches(fileMatches, "main", tmpDir, tmpLicenseFileFi.Name())

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
