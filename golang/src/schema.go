package main

type Schema struct {
	createTableStmt string
	columns         []string
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

	columnName := concatPrefixName(prefix, ast.name)
	if ast.value == JSON_VALUE_OBJECT {
		// Register named objects themselves as a column.
		if ast.name != "" {
			escaped := escapeString(columnName)
			schema.columns = append(schema.columns, escaped)
		}

		// Register the objects members as columns, prefixed with the
		// current column name.
		for i := 0; i < len(ast.members); i++ {
			collectColumns(ast.members[i], schema, columnName)
		}
		return
	}
	escaped := escapeString(columnName)
	schema.columns = append(schema.columns, escaped)
}

func buildTableSchema(ast *ASTNode) Schema {
	schema := Schema{}
	if ast.value == JSON_VALUE_OBJECT {
		collectColumns(ast, &schema, "")
		if len(ast.members) == 0 {
			schema.columns = append(schema.columns, "INTERNAL_PLACEHOLDER")
		}
	} else if ast.value == JSON_VALUE_ARRAY {
		// Use the first entry of the array as the schema.
		// TODO: We should probably add a flag to relax this assumption and take
		//  the set of all values ioIn the input.
		collectColumns(ast.values[0], &schema, "")
		if len(ast.values) == 0 {
			schema.columns = append(schema.columns, "INTERNAL_PLACEHOLDER")
		}
	}

	// Build our 'CREATE TABLE ...' statement from our columns.
	schema.createTableStmt = "CREATE TABLE [] ("
	for i := 0; i < len(schema.columns); i++ {
		schema.createTableStmt += schema.columns[i]
		if i < len(schema.columns)-1 {
			schema.createTableStmt += ","
		}
	}
	schema.createTableStmt += ")"
	return schema
}
