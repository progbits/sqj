package sql

import (
	"sort"
	"testing"
)

func TestExtractIdentifiers_Columns(t *testing.T) {
	// Arrange.
	type TestCase struct {
		statement string
		expected []string
	}
	cases := []TestCase{
		{
			"SELECT a FROM b;",
			[]string{"a"},
		},
		{
			"SELECT a, b, c FROM d WHERE a > 3 AND b < 2;",
			[]string{"a", "b", "c"},
		},
		{
			"SELECT a FROM b WHERE a < 2 AND c > 5 AND d IS NOT 5;",
			[]string{"a", "c", "d"},
		},
		{
			"SELECT a FROM (SELECT b FROM c);",
			[]string{"a", "b"},
		},
	}

	for _, test := range cases {
		// Act
		stmt := parseStatement(test.statement)
		identifiers := ExtractIdentifiers(&stmt, Column)
		sort.Strings(identifiers)

		// Assert.
		if len(identifiers) != len(test.expected) {
			t.Error("unexpected number of identifiers")
		}

		for i := 0; i < len(identifiers); i++ {
			found := false
			for j := 0; j < len(test.expected); j++ {
				if identifiers[i] == test.expected[j] {
					found = true
					break
				}
			}

			if !found {
				t.Error("unexpected identifier")
			}
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
