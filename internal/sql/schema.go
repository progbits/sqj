package sql

import (
	"github.com/progbits/sqjson/internal/util"
	"strings"
)

type SqlSchema struct {
	Columns          [][]string
	CreateTableStmts []string
}

// schemaFromStmt builds a collection of 'CREATE TABLE ...' statements from a
// SQL AST.
func SchemasFromStmt(stmt *SelectStmt) SqlSchema {
	columns := ExtractIdentifiers(stmt, Column)
	tables := ExtractIdentifiers(stmt, Table)

	// This is currently a bit horrible, we should really be returning columns
	// segregated by table.
	orderedColumns := make([][]string, 0)
	createTableStmts := make([]string, 0)
	for i := 0; i < len(tables); i++ {
		tableColumns := make([]string, 0)
		unique := make(map[string]bool)
		stmt := "CREATE TABLE IF NOT EXISTS " + tables[i] + "("
		for j := 0; j < len(columns); j++ {
			parts := strings.Split(columns[j], ".")
			if len(tables) > 1 && len(parts) > 1 && parts[0] != tables[i] {
				continue
			}

			columnName := parts[len(parts)-1]
			if _, ok := unique[columnName]; ok {
				continue
			}

			stmt += util.EscapeString(columnName)
			if j < len(columns)-1 {
				stmt += ","
			}
			tableColumns = append(tableColumns, columnName)
			unique[columnName] = true
		}
		stmt = strings.TrimRight(stmt, ",")
		stmt += ");"

		orderedColumns = append(orderedColumns, tableColumns)
		createTableStmts = append(createTableStmts, stmt)
	}

	return SqlSchema{
		Columns:          orderedColumns,
		CreateTableStmts: createTableStmts,
	}
}
