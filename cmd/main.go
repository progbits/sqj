package main

import (
	"bytes"
	"fmt"
	"github.com/progbits/sqjson/internal/json"
	"github.com/progbits/sqjson/internal/sql"
	"github.com/progbits/sqjson/internal/vtable"
	"github.com/spf13/cobra"

	"io"
	"os"
)

// Global configuration.
var ioIn io.Reader = os.Stdin
var ioOut io.Writer = os.Stdout
var ioErr io.Writer = os.Stderr

type rootCmdVars struct {
	query      string
	inputFiles []string
	nth        string
	compact    bool
}

func runRootCmd(vars *rootCmdVars, cmd *cobra.Command, args []string) {
	var err error

	// Parse the SQL query.
	scanner := sql.NewScanner([]byte(vars.query))
	sqlParser := sql.NewParser(scanner)
	stmt := sqlParser.Parse()

	// Excess arguments after the query string are treated as files and mean we
	// do not read from stdin. A single file named "-" is  treated as an alias
	// for stdin.
	fin := ioIn
	if len(vars.inputFiles) > 0 && vars.inputFiles[0] != "-" {
		fin, err = os.Open(vars.inputFiles[0])
		if err != nil {
			panic("failed to open file")
		}
	}

	// Read the input file to a buffer.
	buf := bytes.NewBuffer(nil)
	_, _ = io.Copy(buf, fin)

	// Tokenize the input data.
	tokenizer := json.Tokenizer{
		Buf: string(buf.Bytes()),
	}
	tokenizer.Tokenize()

	// Parse the token stream.
	jsonParser := json.Parser{
		Tokens: tokenizer.Tokens,
	}
	jsonParser.Parse()

	// Query the virtual table to generate our result ASTs.
	clientData := vtable.ClientData{
		JsonAst: &jsonParser.Ast,
		SqlAst:  &stmt,
		Query:   vars.query,
	}
	result := vtable.Exec(&clientData)

	for _, node := range result {
		json.PrettyPrint(ioOut, node, vars.compact)
		_, _ = fmt.Fprintf(ioOut, "\n")
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "sqj 'QUERY' [FILE]",
		Short: "Query JSON with SQL",
		Long:  `Query JSON with SQL`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var nth string
			var compact bool

			cmd.Flags().StringVarP(&nth, "nth", "n", "", "The nth flag")
			cmd.Flags().BoolVarP(&compact, "compact", "c", false, "The compact flag")

			vars := &rootCmdVars{
				query:      args[0],
				inputFiles: args[1:],
				nth:        nth,
				compact:    compact,
			}
			runRootCmd(vars, cmd, args)
		},
	}

	rootCmd.Execute()
}
