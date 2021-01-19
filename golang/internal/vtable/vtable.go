package vtable

import (
	"database/sql"
	"github.com/progbits/sqjson/internal/json"
	"github.com/progbits/sqjson/internal/util"
	"log"

	"github.com/mattn/go-sqlite3"
)

// Default database driver. This is a global so our integration tests can
// register their own drivers without stepping on each others toes.
var Driver = "sqlite_with_extensions"

type jsonModule struct {
	ast    *json.ASTNode
	query  string
	schema json.Schema
	row    int
}

func (m *jsonModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab(m.schema.CreateTableStmt)
	if err != nil {
		return nil, err
	}
	return &jsonTable{ast: m.ast, schema: m.schema, query: m.query, row: &m.row}, nil
}

func (m *jsonModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *jsonModule) DestroyModule() {}

type jsonTable struct {
	ast    *json.ASTNode
	query  string
	schema json.Schema
	row    *int
}

func (v *jsonTable) Open() (sqlite3.VTabCursor, error) {
	return &jsonCursor{v}, nil
}

func (v *jsonTable) BestIndex(csts []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	return &sqlite3.IndexResult{}, nil
}

func (v *jsonTable) Disconnect() error { return nil }
func (v *jsonTable) Destroy() error    { return nil }

type jsonCursor struct {
	*jsonTable
}

func (vc *jsonCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	columnName := util.UnescapeString(vc.jsonTable.schema.Columns[col])

	var searchNode *json.ASTNode = nil
	if vc.ast.Value == json.JSON_VALUE_OBJECT {
		searchNode = vc.jsonTable.ast
	} else if vc.ast.Value == json.JSON_VALUE_ARRAY {
		searchNode = vc.jsonTable.ast.Values[*vc.row]
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
	if vc.ast.Value == json.JSON_VALUE_OBJECT {
		return *vc.jsonTable.row > 0
	}
	return *vc.jsonTable.row > len(vc.ast.Values)-1
}

func (vc *jsonCursor) Rowid() (int64, error) {
	return int64(*vc.jsonTable.row), nil
}

func (vc *jsonCursor) Close() error {
	return nil
}

func Exec(ast *json.ASTNode, schema json.Schema, query string) []*json.ASTNode {
	jsonModule := jsonModule{
		ast:    ast,
		query:  query,
		schema: schema,
	}

	sql.Register(Driver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("sqjson", &jsonModule)
		},
	})

	db, err := sql.Open(Driver, ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE VIRTUAL TABLE [] USING sqjson")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	nodes := make([]*json.ASTNode, 0)
	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		columns, _ := rows.Columns()
		for i := 0; i < len(columns); i++ {
			var columnNode *json.ASTNode = nil
			if ast.Value == json.JSON_VALUE_OBJECT {
				columnNode = json.FindNode(ast, columns[i])
			} else if ast.Value == json.JSON_VALUE_ARRAY {
				columnNode = json.FindNode(ast.Values[jsonModule.row], columns[i])
			}

			if columnNode != nil {
				nodes = append(nodes, columnNode)
			}
		}
	}
	_ = db.Close()

	return nodes
}
