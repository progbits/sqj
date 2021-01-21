package sql

import "github.com/progbits/sqjson/internal/util"

// schemaFromStmt builds a collection of 'CREATE TABLE ...' statements from a
// SQL AST.
func SchemasFromStmt(stmt *SelectStmt) []string {
	columns := ExtractIdentifiers(stmt, Column)
	tables := ExtractIdentifiers(stmt, Table)

	// This is currently a bit horrible, we should really be returning columns
	// segregated by table.
	createTableStmts := make([]string, 0)
	for i := 0; i < len(tables); i++ {
		stmt := "CREATE TABLE IF NOT EXISTS " + tables[i] + "("
		for j := 0; j < len(columns); j++ {
			stmt += util.EscapeString(columns[j])
			if j < len(columns)-1 {
				stmt += ","
			}
		}
		stmt += ");"
		createTableStmts = append(createTableStmts, stmt)
	}
	return createTableStmts
}
