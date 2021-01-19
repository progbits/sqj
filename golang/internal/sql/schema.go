package sql

type Schema struct {
	columns []string
	tables  []string
}

func extractSchemaImpl(stmt SelectStmt, schema *Schema) {

}

func ExtractSchema(stmt SelectStmt) Schema {
	schema := Schema{}
	extractSchemaImpl(stmt, &schema)
	return schema
}
