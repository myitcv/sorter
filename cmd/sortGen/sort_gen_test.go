package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {
	matches, err := getMatches("_testFiles/", "to_parse.go", "main")

	if err != nil {
		t.Fatalf("Expected err to be nil, got %v", err)
	}

	// number of types matched
	if len(matches) != 2 {
		t.Fatalf("We got %v matches instead of 2", len(matches))
	}

	funs := matches["person"]

	if len(funs) != 2 {
		t.Fatalf("We got %v function matches instead of 2", len(funs))
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
