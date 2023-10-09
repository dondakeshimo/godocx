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
$ go install github.com/dondakeshimo/godocx/cmd/godocx@latest
```

## Usage

```
$ godocx -h
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
```
