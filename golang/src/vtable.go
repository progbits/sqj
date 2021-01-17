package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

// Default database driver. This is a global so our integration tests can
// register their own drivers without stepping on each others toes.
var driver = "sqlite_with_extensions"

type jsonModule struct {
	ast     *ASTNode
	options Options
	schema  Schema
	row     int
}

func (m *jsonModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab(m.schema.createTableStmt)
	if err != nil {
		return nil, err
	}
	return &jsonTable{ast: m.ast, schema: m.schema, options: m.options, row: &m.row}, nil
}

func (m *jsonModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *jsonModule) DestroyModule() {}

type jsonTable struct {
	ast     *ASTNode
	options Options
	schema  Schema
	row     *int
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
	columnName := unescapeString(vc.jsonTable.schema.columns[col])

	var searchNode *ASTNode = nil
	if vc.ast.value == JSON_VALUE_OBJECT {
		searchNode = vc.jsonTable.ast
	} else if vc.ast.value == JSON_VALUE_ARRAY {
		searchNode = vc.jsonTable.ast.values[*vc.row]
	}
	columnNode := findNode(searchNode, columnName)

	if columnNode == nil {
		c.ResultNull()
		return nil
	}

	switch columnNode.value {
	case JSON_VALUE_OBJECT, JSON_VALUE_ARRAY:
		break
	case JSON_VALUE_NUMBER:
		c.ResultDouble(columnNode.number)
	case JSON_VALUE_STRING:
		c.ResultText(columnNode.string)
	case JSON_VALUE_NULL:
		c.ResultNull()
	case JSON_VALUE_TRUE:
		c.ResultBool(true)
	case JSON_VALUE_FALSE:
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
	if vc.ast.value == JSON_VALUE_OBJECT {
		return *vc.jsonTable.row > 0
	}
	return *vc.jsonTable.row > len(vc.ast.values)-1
}

func (vc *jsonCursor) Rowid() (int64, error) {
	return int64(*vc.jsonTable.row), nil
}

func (vc *jsonCursor) Close() error {
	return nil
}

func exec(ast *ASTNode, schema Schema, options Options) {
	jsonModule := jsonModule{
		ast:     ast,
		options: options,
		schema:  schema,
	}

	/*	// Only register "sqlite3_with_extensions" driver if is not already registered.
		drivers := sql.Drivers()
		sqliteRegistered := false
		for _, driver := range drivers {
			if driver == "sqlite3_with_extensions" {
				sqliteRegistered = true
				break;
			}
		}*/

	//if !sqliteRegistered {
	sql.Register(driver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("sqjson", &jsonModule)
		},
	})
	//
	//}

	db, err := sql.Open(driver, ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE VIRTUAL TABLE [] USING sqjson")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(options.query)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		columns, _ := rows.Columns()
		for i := 0; i < len(columns); i++ {
			var columnNode *ASTNode = nil
			if ast.value == JSON_VALUE_OBJECT {
				columnNode = findNode(ast, columns[i])
			} else if ast.value == JSON_VALUE_ARRAY {
				columnNode = findNode(ast.values[jsonModule.row], columns[i])
			}

			if columnNode != nil {
				prettyPrint(ioOut, columnNode, false)
				_, _ = fmt.Fprintf(ioOut, "\n")
			}
		}
	}

	_ = db.Close()
}
