/*
Godocx extracts Go Doc Comment and makes a JSON file.

Usage:

	gofmt [flags] [dirPath ...]

The flags are:

	-o file
		Write result to the file given path instead of stdout.
	-h -help
		Show help.

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
	// NOTE: This usage is not DRY. If you edit this usage, edit the following
	//   - README.md
	//   - package doc
	usage = `
Godocx extracts Go Doc Comment and makes a JSON file.

Usage:

	gofmt [flags] [dirPath ...]

The flags are:

	-o file
		Write result to the file given path instead of stdout.
	-h -help
		Show help.

Godocx supports multiple dirPath which is both absolute and relative path.
Although, it does not support a wild card, or a three dots syntax like gofmt.
Please input each dirPaths separated by a space.
	`
)

// main func is the entry point of the program.
func main() {
	var (
		h    = flag.Bool("h", false, "show help")
		help = flag.Bool("help", false, "show help")
	)
	flag.Parse()

	if *h || *help {
		fmt.Println(usage)
		os.Exit(0)
	}
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
