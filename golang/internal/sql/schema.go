package sql

import (
	"github.com/progbits/sqjson/internal/util"
	"strings"
)

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
			parts := strings.Split(columns[j], ".")
			if len(parts) > 1 && parts[0] != tables[i] {
				continue
			}
			if len(parts) == 1 {
				stmt += util.EscapeString(parts[0])
			} else {
				stmt += util.EscapeString(parts[len(parts)-1])
			}
			if j < len(columns)-1 {
				stmt += ","
			}
		}
		stmt = strings.TrimRight(stmt, ",")
		stmt += ");"
		createTableStmts = append(createTableStmts, stmt)
	}
	return createTableStmts
}
