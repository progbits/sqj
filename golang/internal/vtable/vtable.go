package vtable

import (
	"database/sql"
	"fmt"
	"github.com/progbits/sqjson/internal/json"
	sqlj "github.com/progbits/sqjson/internal/sql"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// Default database driver. This is a global so our integration tests can
// register their own drivers without stepping on each others toes.
var Driver = "sqlite_with_extensions"

type ClientData struct {
	JsonAst *json.ASTNode
	SqlAst  *sqlj.SelectStmt
	Query   string
}

type jsonModule struct {
	clientData      *ClientData
	createTableStmt *string
	result          *map[string]*json.ASTNode
}

func (m *jsonModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	if len(args) < 1 {
		panic("expected table name as argument")
	}

	err := c.DeclareVTab(*m.createTableStmt)
	if err != nil {
		return nil, err
	}

	table := &jsonTable{
		clientData: m.clientData,
		table:      args[2],
		result:     m.result,
	}

	return table, nil
}

func (m *jsonModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *jsonModule) DestroyModule() {}

type jsonTable struct {
	clientData *ClientData
	table      string
	x          int
	y          int
	result     *map[string]*json.ASTNode
}

func (v *jsonTable) Open() (sqlite3.VTabCursor, error) {
	queryRootNode := v.clientData.JsonAst
	var currentNode *json.ASTNode = nil
	if v.table == "[]" {
		// Querying top level node.
		currentNode = v.clientData.JsonAst
	} else {
		// Querying a nested member.
		if v.clientData.JsonAst.Value == json.JSON_VALUE_OBJECT {
			currentNode = json.FindNode(v.clientData.JsonAst, v.table)
			queryRootNode = currentNode
		} else if v.clientData.JsonAst.Value == json.JSON_VALUE_ARRAY {
			currentNode = json.FindNode(v.clientData.JsonAst.Values[0], v.table)
		} else {
			panic("expected an object or an array")
		}
	}

	allColumns := sqlj.ExtractIdentifiers(v.clientData.SqlAst, sqlj.Column)
	var tableColumns []string
	for _, column := range allColumns {
		parts := strings.Split(column, ".")
		if len(parts) > 1 && parts[0] != v.table || len(parts) == 1 {
			continue
		} else {
			tableColumns = append(tableColumns, parts[len(parts)-1])
		}
	}

	if len(tableColumns) == 0 {
		tableColumns = allColumns
	}

	if currentNode == nil {
		panic("unable to locate JSON AST node for table")
	}

	// Construct a new cursor with the column mappings for the current table.
	cursor := &jsonCursor{
		jsonTable: v,
		current:   currentNode,
		queryRoot: queryRootNode,
		columns:   tableColumns,
	}
	return cursor, nil
}

func (v *jsonTable) BestIndex(csts []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	used := make([]bool, len(csts))
	for c, cst := range csts {
		if cst.Usable && cst.Op == sqlite3.OpEQ {
			used[c] = true
		}
	}

	return &sqlite3.IndexResult{
		IdxNum: 0,
		IdxStr: "default",
		Used:   used,
	}, nil
}

func (v *jsonTable) Disconnect() error { return nil }
func (v *jsonTable) Destroy() error    { return nil }

type jsonCursor struct {
	*jsonTable
	current   *json.ASTNode
	queryRoot *json.ASTNode
	columns   []string
	eof       bool
}

func (vc *jsonCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	// Retrieve the original column name.
	columnName := vc.columns[col]
	splitColumnName := strings.Split(columnName, ".")

	// Try and find the AST node corresponding to the column.
	rowNode := vc.current
	var columnNode *json.ASTNode = nil
	if len(splitColumnName) > 1 {
		if rowNode.Value == json.JSON_VALUE_OBJECT {
			columnNode = json.FindNode(rowNode, splitColumnName[0])
		} else if rowNode.Value == json.JSON_VALUE_ARRAY {
			columnNode = json.FindNode(rowNode.Values[vc.y], splitColumnName[0])
		}
	} else {
		if rowNode.Value == json.JSON_VALUE_OBJECT {
			columnNode = json.FindNode(rowNode, columnName)
		} else if rowNode.Value == json.JSON_VALUE_ARRAY {
			columnNode = json.FindNode(rowNode.Values[vc.y], columnName)
		}
	}

	(*vc.result)[columnName] = columnNode

	if columnNode == nil {
		c.ResultNull()
		return nil
	}

	switch columnNode.Value {
	case json.JSON_VALUE_OBJECT, json.JSON_VALUE_ARRAY:
		break
	case json.JSON_VALUE_NUMBER:
		c.ResultDouble(columnNode.Number)
	case json.JSON_VALUE_STRING:
		c.ResultText(columnNode.String)
	case json.JSON_VALUE_NULL:
		c.ResultNull()
	case json.JSON_VALUE_TRUE:
		c.ResultBool(true)
	case json.JSON_VALUE_FALSE:
		c.ResultBool(false)
	}
	return nil
}

func (vc *jsonCursor) Filter(idxNum int, idxStr string, vals []interface{}) error {
	// Reset our cursor.
	if vc.eof {
		vc.x = 0
		vc.y = 0
		vc.current = vc.queryRoot
		vc.eof = false
		return nil
	}

	return nil
}

func (vc *jsonCursor) Next() error {
	// Clear our previous rows results.
	*vc.result = make(map[string]*json.ASTNode)

	// Object queries only execute a single row.
	if vc.queryRoot.Value == json.JSON_VALUE_OBJECT {
		vc.eof = true
		return nil
	}

	// Array queries might be on nested arrays.
	vc.jsonTable.y++
	if vc.jsonTable.y >= len(vc.current.Values) {
		vc.jsonTable.y = 0
		vc.jsonTable.x++
		if vc.jsonTable.x >= len(vc.queryRoot.Values) || vc.queryRoot == vc.current {
			vc.eof = true
			return nil
		}
		vc.current = json.FindNode(vc.queryRoot.Values[vc.jsonTable.x], vc.table)
	}
	return nil
}

func (vc *jsonCursor) EOF() bool {
	return vc.eof
}

func (vc *jsonCursor) Rowid() (int64, error) {
	return int64(vc.jsonTable.x), nil
}

func (vc *jsonCursor) Close() error {
	return nil
}

func Exec(clientData *ClientData) []*json.ASTNode {
	// Extract 'CREATE TABLE ...' statements from SQL AST required to declare
	// the virtual tables for the query.
	createTableStmts := sqlj.SchemasFromStmt(clientData.SqlAst)

	// Register our module and the hook to be invoked on each
	// 'CREATE VIRTUAL TABLE ...' statement.
	result := make(map[string]*json.ASTNode)
	jsonModule := jsonModule{
		clientData: clientData,
		result:     &result,
	}
	sql.Register(Driver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("sqjson", &jsonModule)
		},
	})

	// Open our database connection.
	db, err := sql.Open(Driver, ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// For each table in our query, create the corresponding virtual table.
	// This will call the CreateModule hook to declare the virtual table and
	// initialize the associate jsonTable instance.
	tables := sqlj.ExtractIdentifiers(clientData.SqlAst, sqlj.Table)
	for i := 0; i < len(tables); i++ {
		jsonModule.createTableStmt = &createTableStmts[i]
		_, err = db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE %[1]s USING sqjson", tables[i]))
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err := db.Prepare(clientData.Query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	nodes := make([]*json.ASTNode, 0)
	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		columns, _ := rows.Columns()
		for i := 0; i < len(columns); i++ {
			columnNode := (*jsonModule.result)[columns[i]]
			if columnNode != nil {
				nodes = append(nodes, columnNode)
			}
		}
	}

	return nodes
}
