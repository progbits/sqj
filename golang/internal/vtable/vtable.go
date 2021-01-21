package vtable

import (
	"database/sql"
	"fmt"
	"github.com/progbits/sqjson/internal/json"
	sqlj "github.com/progbits/sqjson/internal/sql"
	"github.com/progbits/sqjson/internal/util"
	"log"

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
	clientData       *ClientData
	createTableStmts []string
	row              int
}

func (m *jsonModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab(m.createTableStmts[0])
	if err != nil {
		return nil, err
	}

	return &jsonTable{clientData: m.clientData, row: &m.row}, nil
}

func (m *jsonModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *jsonModule) DestroyModule() {}

type jsonTable struct {
	clientData *ClientData
	row        *int
}

func (v *jsonTable) Open() (sqlite3.VTabCursor, error) {
	// Construct a new cursor with the column mappings for the current table.
	cursor := &jsonCursor{
		jsonTable: v,
		columns:   sqlj.ExtractIdentifiers(v.clientData.SqlAst, sqlj.Column),
	}
	return cursor, nil
}

func (v *jsonTable) BestIndex(csts []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	return &sqlite3.IndexResult{}, nil
}

func (v *jsonTable) Disconnect() error { return nil }
func (v *jsonTable) Destroy() error    { return nil }

type jsonCursor struct {
	*jsonTable
	columns []string
}

func (vc *jsonCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	// Retrieve the original column name.
	columnName := util.UnescapeString(vc.columns[col])

	// Try and find the AST node corresponding to the column.
	var searchNode *json.ASTNode = nil
	if vc.clientData.JsonAst.Value == json.JSON_VALUE_OBJECT {
		searchNode = vc.jsonTable.clientData.JsonAst
	} else if vc.clientData.JsonAst.Value == json.JSON_VALUE_ARRAY {
		searchNode = vc.jsonTable.clientData.JsonAst.Values[*vc.row]
	}
	columnNode := json.FindNode(searchNode, columnName)

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
	return nil
}

func (vc *jsonCursor) Next() error {
	*vc.jsonTable.row++
	return nil
}

func (vc *jsonCursor) EOF() bool {
	if vc.clientData.JsonAst.Value == json.JSON_VALUE_OBJECT {
		return *vc.jsonTable.row > 0
	}
	return *vc.jsonTable.row > len(vc.clientData.JsonAst.Values)-1
}

func (vc *jsonCursor) Rowid() (int64, error) {
	return int64(*vc.jsonTable.row), nil
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
	jsonModule := jsonModule{
		clientData:       clientData,
		createTableStmts: createTableStmts,
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
		_, err = db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE %s USING sqjson", tables[i]))
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
			var columnNode *json.ASTNode = nil
			if clientData.JsonAst.Value == json.JSON_VALUE_OBJECT {
				columnNode = json.FindNode(clientData.JsonAst, columns[i])
			} else if clientData.JsonAst.Value == json.JSON_VALUE_ARRAY {
				columnNode = json.FindNode(clientData.JsonAst.Values[jsonModule.row], columns[i])
			}

			if columnNode != nil {
				nodes = append(nodes, columnNode)
			}
		}
	}

	return nodes
}
