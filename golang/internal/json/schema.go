package json

import "github.com/progbits/sqjson/internal/util"

type Schema struct {
	CreateTableStmt string
	Columns         []string
}

// Concatenate a prefix an a member name.
//
// JSON objects can either be top level nodes or themselves object members. For
// nested members, column names are a concatenation of the member names,
// separated by $. Top-level objects have no name, so their concatenated name is
// simply the prefix itself with no trailing $.
//
func concatPrefixName(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "$" + name
}

func collectColumns(ast *ASTNode, schema *Schema, prefix string) {
	if ast == nil {
		return
	}

	columnName := concatPrefixName(prefix, ast.Name)
	if ast.Value == JSON_VALUE_OBJECT {
		// Register named objects themselves as a column.
		if ast.Name != "" {
			escaped := util.EscapeString(columnName)
			schema.Columns = append(schema.Columns, escaped)
		}

		// Register the objects members as columns, prefixed with the
		// current column name.
		for i := 0; i < len(ast.Members); i++ {
			collectColumns(ast.Members[i], schema, columnName)
		}
		return
	}
	escaped := util.EscapeString(columnName)
	schema.Columns = append(schema.Columns, escaped)
}
