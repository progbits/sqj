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

func BuildTableSchema(ast *ASTNode) Schema {
	schema := Schema{}
	if ast.Value == JSON_VALUE_OBJECT {
		collectColumns(ast, &schema, "")
		if len(ast.Members) == 0 {
			schema.Columns = append(schema.Columns, "INTERNAL_PLACEHOLDER")
		}
	} else if ast.Value == JSON_VALUE_ARRAY {
		// Use the first entry of the array as the schema.
		// TODO: We should probably add a flag to relax this assumption and take
		//  the set of all values ioIn the input.
		collectColumns(ast.Values[0], &schema, "")
		if len(ast.Values) == 0 {
			schema.Columns = append(schema.Columns, "INTERNAL_PLACEHOLDER")
		}
	}

	// Build our 'CREATE TABLE ...' statement from our columns.
	schema.CreateTableStmt = "CREATE TABLE [] ("
	for i := 0; i < len(schema.Columns); i++ {
		schema.CreateTableStmt += schema.Columns[i]
		if i < len(schema.Columns)-1 {
			schema.CreateTableStmt += ","
		}
	}
	schema.CreateTableStmt += ")"
	return schema
}
