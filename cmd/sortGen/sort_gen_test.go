package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {
	matches, err := getMatches("_testFiles/", "to_parse.go", "main")

	if err != nil {
		t.Fatalf("Expected err to be nil, got %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("We got %v matches instead of 1", len(matches))
	}

	funs := matches["person"]

	if len(funs) != 2 {
		t.Fatalf("We got %v function matches instead of 2", len(funs))
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	err = genMatches(matches, "main", tmpDir)

	if err != nil {
		t.Fatalf("Expected gen err to be nil, got %v", err)
	}

	outFile, err := ioutil.ReadFile(filepath.Join(tmpDir, "gen_person_sorter.go"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(outFile))
}
