// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main // import "myitcv.io/immutable/cmd/immutableGen"

import (
	"flag"
	"fmt"
	"log"
	"os"

	"myitcv.io/gogenerate"
)

const (
	immutableGenCmd = "immutableGen"
)

var (
	fGoGenCmds   gogenCmds
	fLicenseFile = gogenerate.LicenseFileFlag()
	fGoGenLog    = gogenerate.LogFlag()
	fDebug       = flag.Bool("debug", false, "print debug messages")
)

const (
	debug = false
)

func init() {
	flag.Var(&fGoGenCmds, "G", "Path to search for imports (flag can be used multiple times)")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix(immutableGenCmd + ": ")

	gogenerate.DefaultLogLevel(fGoGenLog, gogenerate.LogFatal)

	envFile, ok := os.LookupEnv(gogenerate.GOFILE)
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

	dirFiles, err := gogenerate.FilesContainingCmd(wd, immutableGenCmd)
	if err != nil {
		fatalf("could not determine if we are the first file: %v", err)
	}

	if dirFiles == nil {
		fatalf("cannot find any files containing the %v directive", immutableGenCmd)
	}

	if dirFiles[envFile] != 1 {
		fatalf("expected a single occurrence of %v directive in %v. Got: %v", immutableGenCmd, envFile, dirFiles)
	}

	licenseHeader, err := gogenerate.CommentLicenseHeader(fLicenseFile)
	if err != nil {
		fatalf("could not comment license file: %v", err)
	}

	execute(wd, envPkgName, licenseHeader, fGoGenCmds)
}

func debugf(format string, args ...interface{}) {
	if debug || *fDebug {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func infoln(args ...interface{}) {
	if *fGoGenLog == string(gogenerate.LogInfo) {
		log.Println(args...)
	}
}

func infof(format string, args ...interface{}) {
	if *fGoGenLog == string(gogenerate.LogInfo) {
		log.Printf(format, args...)
	}
}
