# godocx

godocx is a command line tool to extract Go Doc Comments and transform to JSON.
You can currate Go Doc Comments from various perspectives with output from godocx.

`go doc` command provides a way to transform Go Doc Comments to plain text, but it is not easy to currate the output.
`godoc` command build a web server and provides a web interface to view Go Doc Comments, but it is not easy to filter or search.
`go/doc` package provides a way to extract Go Doc Comments, but it is not easy to write codes for each perspectives.

## Features

- Extract Go Doc Comments from multiple packages
- Annotate Go Doc Comments with a prefix `@`

## Installation

```
$ go install github.com/dondakeshimo/godocx@latest
```

## Usage

(WIP)

```
$ godocx -h
Usage of godocx:
  -all
        include all packages
  -c string
        comment type (default "all")
  -d string
        output directory (default "docs")
  -f string
        output format (default "json")
  -p string
        package path (default ".")
  -v    verbose output
```
