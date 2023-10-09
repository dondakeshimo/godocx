/*
Godocx extracts Go Doc Comment and makes a JSON file.

Usage:

	gofmt [flags] [dirPath ...]

The flags are:

	-o file
		Write result to the file given path instead of stdout.

Godocx supports multiple dirPath which is both absolute and relative path.
Although, it does not support a wild card, or a three dots syntax like gofmt.
Please input each dirPaths separated by a space.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/dondakeshimo/godocx"
)

const (
	usage = `
Godocx extracts Go Doc Comment and makes a JSON file.

Usage:

	gofmt [flags] [dirPath ...]

The flags are:

	-o file
		Write result to the file given path instead of stdout.

Godocx supports multiple dirPath which is both absolute and relative path.
Although, it does not support a wild card, or a three dots syntax like gofmt.
Please input each dirPaths separated by a space.
	`
)

// main func is the entry point of the program.
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("godocx: at least one argument required.")
		fmt.Println(usage)
		os.Exit(1)
	}

	pkg, err := godocx.New(args)
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
