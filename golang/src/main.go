package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Options struct {
	compact bool
	nth     int

	query string
}

func usage() {
	_, _ = fmt.Fprintf(ioErr,
		`Usage: sqj [OPTION]... <SQL> [FILE]...
Query JSON with SQL.

    --help    Display this message and exit
    --compact Format output without any extraneous whitespace`)
}

// Global configuration.
var ioIn io.Reader = os.Stdin
var ioOut io.Writer = os.Stdout
var ioErr io.Writer = os.Stderr

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// Parse command line options.
	i := 0
	var err error
	options := Options{nth: -1}
	for i = 1; i < len(os.Args); i++ {
		z := os.Args[i]
		if z[0] != '-' {
			break
		}
		if z[1] == '-' {
			z = z[1:]
		}

		if z == "-help" {
			usage()
		} else if z == "-compact" {
			options.compact = true
		} else if z == "-nth" {
			options.nth, _ = strconv.Atoi(os.Args[i])
		}
	}

	// The query string should be the first argument after options.
	options.query = os.Args[i]
	i++

	// Excess arguments after the query string are treated as files and mean we
	// do not read from stdin. A single file named "-" is  treated as an alias
	// for stdin.
	fin := ioIn
	if i < len(os.Args) && os.Args[i] != "-" {
		fin, err = os.Open(os.Args[i])
		if err != nil {
			panic("failed to open file")
		}
	}

	// Read the input file to a buffer.
	buf := bytes.NewBuffer(nil)
	_, _ = io.Copy(buf, fin)

	// Tokenize the input data.
	tokenizer := Tokenizer{
		buf: string(buf.Bytes()),
	}
	tokenizer.tokenize()

	// Parse the token stream.
	parser := Parser{
		tokens: tokenizer.tokens,
	}
	parser.parse()

	schema := buildTableSchema(&parser.ast)
	exec(&parser.ast, schema, options)
}
