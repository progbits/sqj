package sql

import "testing"

func TestExtractIdentifiers_Columns(t *testing.T) {
	// Arrange.

	// SELECT a FROM (SELECT b FROM c WHERE b == 5) AS d;
	stmt := SelectStmt{
		tableList: TableList{
			source: SelectStmt{
				resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "b"}}},
				tableList:    TableList{source: &IdentifierExpr{value: "c", kind: Table}},
				whereClause: &BinaryExpr{
					operator: EQ,
					left:     &IdentifierExpr{value: "b"},
					right:    &LiteralExpr{value: "5"},
				},
			},
		},
	}

	// Act.
	identifiers := ExtractIdentifiers(&stmt, Column)

	// Assert.
	expectedIdentifiers := []string{
		"a", "b",
	}

	if len(identifiers) != len(expectedIdentifiers) {
		t.Error("unexpected number of identifiers")
	}

	for i := 0; i < len(identifiers); i++ {
		found := false
		for j := 0; j < len(expectedIdentifiers); j++ {
			if identifiers[i] != expectedIdentifiers[j] {
				found = true
				break
			}
		}

		if !found {
			t.Error("unexpected identifier")
		}
	}
}

func TestExtractIdentifiers_Tables(t *testing.T) {
	// Arrange.

	// SELECT * FROM [];
	stmt := SelectStmt{
		resultColumn: []ResultColumn{
			{expr: &StarExpr{}},
		},
		tableList: TableList{
			source: &IdentifierExpr{kind: Table, value: "*"},
		},
	}

	// Act.
	identifiers := ExtractIdentifiers(&stmt, Table)

	// Assert.
	expectedIdentifiers := []string{
		"[]",
	}

	if len(identifiers) != len(expectedIdentifiers) {
		t.Error("unexpected number of identifiers")
	}

	for i := 0; i < len(identifiers); i++ {
		found := false
		for j := 0; j < len(expectedIdentifiers); j++ {
			if identifiers[i] != expectedIdentifiers[j] {
				found = true
				break
			}
		}

		if !found {
			t.Error("unexpected identifier")
		}
	}
}
