// Package godocx is the library for extracting Go documentation.
package godocx

import (
	"errors"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
)

// Value is the documentation for a var or const declaration.
type Value struct {
	Name        string   `json:"name"`
	Doc         string   `json:"doc"`
	Annotations []string `json:"annotations"`
}

// Type is the documentation for a type declaration.
type Type struct {
	Name        string   `json:"name"`
	Doc         string   `json:"doc"`
	Annotations []string `json:"annotations"`
	Consts      []*Value `json:"consts"`
	Vars        []*Value `json:"vars"`
	Funcs       []*Func  `json:"funcs"`
}

// Func is the documentation for a func declaration.
type Func struct {
	Name        string   `json:"name"`
	Doc         string   `json:"doc"`
	Recv        string   `json:"recv"`
	Orig        string   `json:"orig"`
	Level       int      `json:"level"`
	Annotations []string `json:"annotations"`
}

// A Note represents a marked comment starting with "MARKER(uid): note body".
// Any note with a marker of 2 or more upper case [A-Z] letters and a uid of
// at least one character is recognized. The ":" following the uid is optional.
type Note struct {
	UID  string `json:"uid"`
	Body string `json:"body"`
}

// Package is the documentation for an entire package.
type Package struct {
	Name       string             `json:"name"`
	ImportPath string             `json:"importPath"`
	Notes      map[string][]*Note `json:"notes"`
	Consts     []*Value           `json:"consts"`
	Vars       []*Value           `json:"vars"`
	Funcs      []*Func            `json:"funcs"`
	Types      []*Type            `json:"types"`
}

// New creates a new Package from the given package name and import path.
func New(dirPath string) (*Package, error) {
	// To search target files, use build.ImportDir.
	buildPkg, err := build.ImportDir(dirPath, build.ImportComment)
	if err != nil {
		return nil, err
	}

	// include tells parser.ParseDir which files to include.
	// That means the file must be in the build package's GoFiles or CgoFiles
	// list only (no tag-ignored files, tests, swig or other non-Go files).
	include := func(info fs.FileInfo) bool {
		for _, name := range buildPkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range buildPkg.CgoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}

	// analyze a package with ast.
	fset := token.NewFileSet()
	pkgMap, firstErr := parser.ParseDir(fset, dirPath, include, parser.ParseComments)
	if firstErr != nil {
		return nil, firstErr
	}

	// check the number of packages.
	// ParseDir is given a func `include`, so the number of packages should be 1.
	if len(pkgMap) == 0 {
		return nil, errors.New("godocx: no packages found")
	} else if len(pkgMap) > 1 {
		return nil, errors.New("godocx: multiple packages found")
	}

	// analyze a package with doc.
	d := doc.New(pkgMap[buildPkg.Name], ".", doc.AllDecls|doc.AllMethods)

	// convert doc.Package to Package.

	consts := make([]*Value, 0, len(d.Consts))
	for _, v := range d.Consts {
		consts = append(consts, &Value{
			Name:        v.Names[0],
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
		})
	}

	vars := make([]*Value, 0, len(d.Vars))
	for _, v := range d.Vars {
		vars = append(vars, &Value{
			Name:        v.Names[0],
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
		})
	}

	Funcs := make([]*Func, 0, len(d.Funcs))
	for _, v := range d.Funcs {
		Funcs = append(Funcs, &Func{
			Name:        v.Name,
			Doc:         v.Doc,
			Recv:        v.Recv,
			Orig:        v.Orig,
			Level:       v.Level,
			Annotations: extractAnnotations(v.Doc),
		})
	}

	types := make([]*Type, 0, len(d.Types))
	for _, v := range d.Types {
		cs := make([]*Value, 0, len(v.Consts))
		for _, v := range v.Consts {
			cs = append(cs, &Value{
				Name:        v.Names[0],
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
			})
		}
		vs := make([]*Value, 0, len(v.Vars))
		for _, v := range v.Vars {
			vs = append(vs, &Value{
				Name:        v.Names[0],
				Doc:         v.Doc,
				Annotations: extractAnnotations(v.Doc),
			})
		}
		fs := make([]*Func, 0, len(v.Funcs))
		for _, v := range v.Funcs {
			fs = append(fs, &Func{
				Name:        v.Name,
				Doc:         v.Doc,
				Recv:        v.Recv,
				Orig:        v.Orig,
				Level:       v.Level,
				Annotations: extractAnnotations(v.Doc),
			})
		}
		types = append(types, &Type{
			Name:        v.Name,
			Doc:         v.Doc,
			Annotations: extractAnnotations(v.Doc),
			Consts:      cs,
			Vars:        vs,
			Funcs:       fs,
		})
	}

	notes := make(map[string][]*Note, len(d.Notes))
	for k, v := range d.Notes {
		n := make([]*Note, 0, len(v))
		for _, vv := range v {
			n = append(n, &Note{
				UID:  vv.UID,
				Body: vv.Body,
			})
		}
		notes[k] = n
	}

	return &Package{
		Name:       d.Name,
		ImportPath: d.ImportPath,
		Notes:      notes,
		Consts:     consts,
		Vars:       vars,
		Funcs:      Funcs,
		Types:      types,
	}, nil
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
