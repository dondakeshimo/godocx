// Package main is the entry point of the program.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dondakeshimo/godocx"
)

// main func is the entry point of the program.
func main() {
	if len(os.Args) == 1 {
		fmt.Println("Hello, World!")
		os.Exit(0)
	}

	// TODO: Be able to search packages only package's name specified.
	//   - scan scope can be GOROOT/src, GOPATH/src, or current directory.
	/*
		goroot := os.Getenv("GOROOT")
		fmt.Printf("GOROOT: %s\n", goroot)

		gopath := os.Getenv("GOPATH")
		fmt.Printf("GOPATH: %s\n", gopath)

		importPath := os.Args[1]
		buildPkg, err := build.ImportDir(importPath, build.ImportComment)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", buildPkg)
	*/

	// Get the directory name from the command line.
	dir := os.Args[1]

	pkg, err := godocx.New([]string{dir})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := json.Marshal(pkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}
