// Package main is the entry point of the program.
package main

import (
	"encoding/json"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

// Node represents a node of the Go Doc Comment.
type Node struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Doc         string   `json:"doc"`
	Comment     string   `json:"comment"`
	Annotations []string `json:"annotations"`
	Notes       []string `json:"notes"`
}

// Node represents a type node of the Go Doc Comment.
type TypeNode struct {
	Node
	Consts []*Node `json:"consts"`
	Vars   []*Node `json:"vars"`
	Funcs  []*Node `json:"funcs"`
}

// Package represents a package of the Go Doc Comment.
type Package struct {
	Name       string      `json:"name"`
	ImportPath string      `json:"importPath"`
	Consts     []*Node     `json:"consts"`
	Vars       []*Node     `json:"vars"`
	Funcs      []*Node     `json:"funcs"`
	Types      []*TypeNode `json:"types"`
}

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

	fset := token.NewFileSet()
	pkgMap, firstErr := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if firstErr != nil {
		fmt.Println(firstErr)
		os.Exit(1)
	}

	d := doc.New(pkgMap["btime"], ".", doc.AllDecls|doc.AllMethods)

	consts := make([]*Node, 0, len(d.Consts))
	for _, v := range d.Consts {
		consts = append(consts, &Node{
			Name:        v.Names[0],
			Kind:        "const",
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
			Notes:       extractNotes(v.Doc),
		})
	}

	vars := make([]*Node, 0, len(d.Vars))
	for _, v := range d.Vars {
		vars = append(vars, &Node{
			Name:        v.Names[0],
			Kind:        "var",
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
			Notes:       extractNotes(v.Doc),
		})
	}

	Funcs := make([]*Node, 0, len(d.Funcs))
	for _, v := range d.Funcs {
		Funcs = append(Funcs, &Node{
			Name:        v.Name,
			Kind:        "func",
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
			Notes:       extractNotes(v.Doc),
		})
	}

	types := make([]*TypeNode, 0, len(d.Types))
	for _, v := range d.Types {
		cs := make([]*Node, 0, len(v.Consts))
		for _, v := range v.Consts {
			cs = append(cs, &Node{
				Name:        v.Names[0],
				Kind:        "const",
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
				Notes:       extractNotes(v.Doc),
			})
		}
		vs := make([]*Node, 0, len(v.Vars))
		for _, v := range v.Vars {
			vs = append(vs, &Node{
				Name:        v.Names[0],
				Kind:        "var",
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
				Notes:       extractNotes(v.Doc),
			})
		}
		fs := make([]*Node, 0, len(v.Funcs))
		for _, v := range v.Funcs {
			fs = append(fs, &Node{
				Name:        v.Name,
				Kind:        "func",
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
				Notes:       extractNotes(v.Doc),
			})
		}
		types = append(types, &TypeNode{
			Node: Node{
				Name:        v.Name,
				Kind:        "type",
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
				Notes:       extractNotes(v.Doc),
			},
			Consts: cs,
			Vars:   vs,
			Funcs:  fs,
		})
	}

	pkg := &Package{
		Name:       d.Name,
		ImportPath: d.ImportPath,
		Consts:     consts,
		Vars:       vars,
		Funcs:      Funcs,
		Types:      types,
	}

	b, err := json.Marshal(pkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}

var annotationRegExp = regexp.MustCompile(`@(\w+)`)

// extractAnnotations extracts annotations from the comment.
func extractAnnotations(comment string) []string {
	annotations := annotationRegExp.FindAllString(comment, -1)
	if annotations == nil {
		return []string{}
	}

	return annotations
}

// extractNotes extracts notes from the comment.
func extractNotes(comment string) []string {
	notes := make([]string, 0)
	if strings.Contains(comment, "TODO") {
		notes = append(notes, "TODO")
	}
	if strings.Contains(comment, "FIXME") {
		notes = append(notes, "FIXME")
	}
	if strings.Contains(comment, "XXX") {
		notes = append(notes, "XXX")
	}
	if strings.Contains(comment, "BUG") {
		notes = append(notes, "BUG")
	}
	if strings.Contains(comment, "NOTE") {
		notes = append(notes, "NOTE")
	}
	if strings.Contains(comment, "HACK") {
		notes = append(notes, "HACK")
	}
	if strings.Contains(comment, "OPTIMIZE") {
		notes = append(notes, "OPTIMIZE")
	}
	if strings.Contains(comment, "WARNING") {
		notes = append(notes, "WARNING")
	}
	if strings.Contains(comment, "ERROR") {
		notes = append(notes, "ERROR")
	}
	return notes
}
