package sql

import (
	"testing"
)

// parseStatement is a helper method to parse a statement for testing.
func parseStatement(stmt string) SelectStmt {
	scanner := NewScanner([]byte(stmt))
	parser := Parser{
		scanner: scanner,
	}
	return parser.Parse()
}

func TestParseExpr(t *testing.T) {
	type TestCase struct {
		statement string
		expected  []Expr
	}

	var cases = [...]TestCase{
		// literal values
		{"SELECT 5;", []Expr{&LiteralExpr{value: "5"}}},
		//{"SELECT \"hello, world\";", []Expr{}{LiteralExpr{value: "hello, world"}}},
		{"SELECT a;", []Expr{&IdentifierExpr{value: "a"}}},
		{"SELECT α, γ, δ, ϵ, ζ, θ, μ, ψ;", []Expr{
			&IdentifierExpr{value: "α"},
			&IdentifierExpr{value: "γ"},
			&IdentifierExpr{value: "δ"},
			&IdentifierExpr{value: "ϵ"},
			&IdentifierExpr{value: "ζ"},
			&IdentifierExpr{value: "θ"},
			&IdentifierExpr{value: "μ"},
			&IdentifierExpr{value: "ψ"},
		}},
		// unary operators
		{"SELECT -a;", []Expr{&UnaryExpr{operator: MINUS, expr: &IdentifierExpr{value: "a"}}}},
		{"SELECT -(a + b > 2);", []Expr{
			&UnaryExpr{
				operator: MINUS,
				expr: &BinaryExpr{
					operator: GT,
					left: &BinaryExpr{
						operator: PLUS,
						left:     &IdentifierExpr{value: "a"},
						right:    &IdentifierExpr{value: "b"},
					},
					right: &LiteralExpr{value: "2"},
				},
			},
		}},
		{"SELECT +a;", []Expr{&UnaryExpr{operator: PLUS, expr: &IdentifierExpr{value: "a"}}}},
		{"SELECT ~a;", []Expr{&UnaryExpr{operator: BITNOT, expr: &IdentifierExpr{value: "a"}}}},
		{"SELECT ~(a + b > 2);", []Expr{
			&UnaryExpr{
				operator: BITNOT,
				expr: &BinaryExpr{
					operator: GT,
					left: &BinaryExpr{
						operator: PLUS,
						left:     &IdentifierExpr{value: "a"},
						right:    &IdentifierExpr{value: "b"},
					},
					right: &LiteralExpr{value: "2"},
				},
			},
		}},
		{"SELECT NOT a;", []Expr{&UnaryExpr{operator: NOT, expr: &IdentifierExpr{value: "a"}}}},
		{"SELECT NOT(a + b > 2);", []Expr{
			&UnaryExpr{
				operator: NOT,
				expr: &BinaryExpr{
					operator: GT,
					left: &BinaryExpr{
						operator: PLUS,
						left:     &IdentifierExpr{value: "a"},
						right:    &IdentifierExpr{value: "b"},
					},
					right: &LiteralExpr{value: "2"},
				},
			},
		}},
		// binary operators, precedence order
		{"SELECT a || b;", []Expr{&BinaryExpr{operator: CONCAT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a * b;", []Expr{&BinaryExpr{operator: STAR, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a / b;", []Expr{&BinaryExpr{operator: SLASH, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a % b;", []Expr{&BinaryExpr{operator: REM, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a + b;", []Expr{&BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a - b;", []Expr{&BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a << b;", []Expr{&BinaryExpr{operator: LSHIFT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a >> b;", []Expr{&BinaryExpr{operator: RSHIFT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a & b;", []Expr{&BinaryExpr{operator: BITAND, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a | b;", []Expr{&BinaryExpr{operator: BITOR, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a < b;", []Expr{&BinaryExpr{operator: LT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a <= b;", []Expr{&BinaryExpr{operator: LE, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a > b;", []Expr{&BinaryExpr{operator: GT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a >= b;", []Expr{&BinaryExpr{operator: GE, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a = b;", []Expr{&BinaryExpr{operator: EQ, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a == b;", []Expr{&BinaryExpr{operator: EQ, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a != b;", []Expr{&BinaryExpr{operator: NE, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a <> b;", []Expr{&BinaryExpr{operator: NE, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a NOT b;", []Expr{&BinaryExpr{operator: NOT, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a IN b;", []Expr{&BinaryExpr{operator: IN, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		//	{"SELECT a LIKE b;", []Expr{&BinaryExpr{operator: LIKE, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}}}},
		//	{"SELECT a GLOB b;", []Expr{&BinaryExpr{operator: GLOB, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}}}},
		//	{"SELECT a MATCH b;", []Expr{&BinaryExpr{operator: MATCH, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}}}},
		//	{"SELECT a REGEXP b;", []Expr{&BinaryExpr{operator: REGEXP, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}}}},
		{"SELECT a AND b;", []Expr{&BinaryExpr{operator: AND, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a OR b;", []Expr{&BinaryExpr{operator: OR, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		// TODO: function call expression
		// parenthesised expressions
		{"SELECT (a)", []Expr{&IdentifierExpr{value: "a"}}},
		{"SELECT (((((a)))))", []Expr{&IdentifierExpr{value: "a"}}},
		{"SELECT ((a + 5) * 3)", []Expr{
			&BinaryExpr{
				operator: STAR,
				left: &BinaryExpr{
					operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &LiteralExpr{value: "5"},
				},
				right: &LiteralExpr{value: "3"},
			},
		}},
		// CAST expressions
		{"SELECT CAST (a + 5 AS value);", []Expr{
			&CastExpr{
				typeName: "value",
				expr:     &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &LiteralExpr{value: "5"}},
			},
		}},
		// COLLATE expressions
		{"SELECT a + 5 COLLATE BINARY;", []Expr{
			&CollateExpr{
				collationName: "BINARY",
				expr: &BinaryExpr{
					operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &LiteralExpr{value: "5"},
				},
			},
		}},
		// string match (LIKE, GLOB, REGEXP, MATCH) expressions
		// TODO: Handle ESCAPE expression
		{"SELECT a LIKE b;", []Expr{&StringMatchExpr{operator: LIKE, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a GLOB b;", []Expr{&StringMatchExpr{operator: GLOB, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a REGEXP b;", []Expr{&StringMatchExpr{operator: REGEXP, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a MATCH b;", []Expr{&StringMatchExpr{operator: MATCH, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		//	{"SELECT (a + b) MATCH (b - d) ESCAPE (c * e);", []Expr{}{
		//		&StringMatchExpr{
		//			operator:   MATCH,
		//			left:       BinaryExpr{operator: PLUS, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}},
		//			right:      BinaryExpr{operator: MINUS, left: IdentifierExpr{value: "b"}, right: IdentifierExpr{value: "d"}},
		//			escapeExpr: BinaryExpr{operator: STAR, left: IdentifierExpr{value: "c"}, right: IdentifierExpr{value: "e"}},
		//		},
		//	}},
		// inverse string match
		{"SELECT a NOT LIKE b;", []Expr{&StringMatchExpr{operator: LIKE, inverse: true, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a NOT GLOB b;", []Expr{&StringMatchExpr{operator: GLOB, inverse: true, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a NOT REGEXP b;", []Expr{&StringMatchExpr{operator: REGEXP, inverse: true, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a NOT MATCH b;", []Expr{&StringMatchExpr{operator: MATCH, inverse: true, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT (a + b) NOT MATCH (b - d);", []Expr{
			&StringMatchExpr{
				operator: MATCH,
				inverse:  true,
				left:     &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
				right:    &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
			},
		}},
		// TODO: ISNULL, NOTNULL & NOT NULL
		// IS expressions
		{"SELECT a IS b;", []Expr{&IsExpr{left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT a IS NOT b;", []Expr{&IsExpr{inverse: true, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}}}},
		{"SELECT (a + b) IS NOT (b - d);", []Expr{
			&IsExpr{
				inverse: true,
				left:    &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
				right:   &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
			},
		}},
		// between expressions
		{"SELECT a BETWEEN b AND c;", []Expr{&BetweenExpr{
			inverse: false,
			expr:    &IdentifierExpr{value: "a"},
			left:    &IdentifierExpr{value: "b"},
			right:   &IdentifierExpr{value: "c"},
		}}},
		{"SELECT (a + b) BETWEEN (b - c) AND (b AND d);", []Expr{&BetweenExpr{
			inverse: false,
			expr:    &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
			left:    &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "c"}},
			right:   &BinaryExpr{operator: AND, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
		}}},
		{"SELECT a NOT BETWEEN b AND c;", []Expr{&BetweenExpr{
			inverse: true,
			expr:    &IdentifierExpr{value: "a"},
			left:    &IdentifierExpr{value: "b"},
			right:   &IdentifierExpr{value: "c"},
		}}},
		{"SELECT (a + b) NOT BETWEEN (b - c) AND (b AND d);", []Expr{&BetweenExpr{
			inverse: true,
			expr:    &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
			left:    &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "c"}},
			right:   &BinaryExpr{operator: AND, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
		}}},
		// TODO: IN expression
		// case expression
		{"SELECT CASE WHEN b > c THEN d WHEN e < f THEN G ELSE a END;", []Expr{&CaseExpr{
			expr: nil,
			when: []Expr{
				&BinaryExpr{operator: GT, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "c"}},
				&BinaryExpr{operator: LT, left: &IdentifierExpr{value: "e"}, right: &IdentifierExpr{value: "f"}},
			},
			then: []Expr{
				&IdentifierExpr{value: "d"},
				&IdentifierExpr{value: "G"},
			},
			elseExpr: &IdentifierExpr{value: "a"},
		}}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if len(stmt.resultColumn) != len(_case.expected) {
			t.Error("unexpected number of expressions")
		}
		for i := 0; i < len(stmt.resultColumn); i++ {
			if !eqExpr(stmt.resultColumn[i].expr, _case.expected[i]) {
				t.Error("unexpected expression")
			}
		}
	}
}

func TestParseInExpr(t *testing.T) {
	/*	// Arrange
		stmt := "SELECT EXISTS(SELECT a FROM b WHERE a == 3);"
		expected := SelectStmt{
			resultColumn: []ResultColumn{{
				expr: &ExistsExpr{
					selectStmt: SelectStmt{
						resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "a"}}},
						fromClause:   TableExpr{source: &IdentifierExpr{value: "b", kind: Table}},
						whereClause: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a"},
							right:    &LiteralExpr{value: "3"},
						},
					},
				},
			}},
		}

		// Act / Assert
		parsedStmt := parseStatement(stmt)
		checkStatement(t, parsedStmt, expected)*/
}

func TestParseInverseInExpr(t *testing.T) {
	/*	// Arrange
		stmt := "SELECT NOT EXISTS(SELECT a FROM b WHERE a == 3);"
		expected := SelectStmt{
			resultColumn: []ResultColumn{{
				expr: &ExistsExpr{
					inverse: true,
					selectStmt: SelectStmt{
						resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "a"}}},
						fromClause:   []TableExpr{source: &IdentifierExpr{value: "b", kind: Table}},
						whereClause: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a"},
							right:    &LiteralExpr{value: "3"},
						},
					},
				},
			}},
		}

		// Act / Assert
		parsedStmt := parseStatement(stmt)
		checkStatement(t, parsedStmt, expected)*/
}

func TestParseDistinctAll(t *testing.T) {
	type TestCase struct {
		statement string
		expected  SelectStmt
	}

	var cases = [...]TestCase{
		{"SELECT DISTINCT *;", SelectStmt{isDistinct: true}},
		{"SELECT ALL *;", SelectStmt{isAll: true}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if stmt.isDistinct != _case.expected.isDistinct {
			t.Fatalf("unexpected DISTINCT")
		}
		if stmt.isAll != _case.expected.isAll {
			t.Fatalf("unexpected ALL")
		}
	}
}

func TestParseResultColumns(t *testing.T) {
	type TestCase struct {
		statement string
		expected  []ResultColumn
	}

	var cases = [...]TestCase{
		{"SELECT *;", []ResultColumn{{expr: &StarExpr{}}}},
		{"SELECT a;", []ResultColumn{{expr: &IdentifierExpr{value: "a"}}}},
		{"SELECT a, b, c;", []ResultColumn{
			{expr: &IdentifierExpr{value: "a"}},
			{expr: &IdentifierExpr{value: "b"}},
			{expr: &IdentifierExpr{value: "c"}},
		}},
		{"SELECT test_table.*;", []ResultColumn{{expr: &IdentifierExpr{value: "test_table."}}}},
		{"SELECT test_column_a AS a, test_column_b AS b;", []ResultColumn{
			{expr: &IdentifierExpr{value: "test_column_a"}, alias: "a"},
			{expr: &IdentifierExpr{value: "test_column_b"}, alias: "b"},
		}},
		{"SELECT (a + b - c * d NOT z) AS a;", []ResultColumn{{expr: &BinaryExpr{
			operator: NOT,
			left: &BinaryExpr{
				operator: MINUS,
				left: &BinaryExpr{
					operator: PLUS,
					left:     &IdentifierExpr{value: "a"},
					right:    &IdentifierExpr{value: "b"},
				},
				right: &BinaryExpr{
					operator: STAR,
					left:     &IdentifierExpr{value: "c"},
					right:    &IdentifierExpr{value: "d"},
				},
			},
			right: &IdentifierExpr{value: "z"},
		},
			alias: "a",
		}}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if len(stmt.resultColumn) != len(_case.expected) {
			t.Error("unexpected number of expressions")
		}

		for i := 0; i < len(stmt.resultColumn); i++ {
			if !eqExpr(stmt.resultColumn[i].expr, _case.expected[i].expr) {
				t.Error("unexpected expression")
			}
		}
	}
}

func TestParseTableList(t *testing.T) {
	type TestCase struct {
		statement string
		expected  []JoinedTable
	}

	var cases = [...]TestCase{
		{"SELECT * FROM a;", []JoinedTable{
			{source: &IdentifierExpr{value: "a", kind: Table}},
		}},
		{"SELECT * FROM a, b, c;", []JoinedTable{
			{source: &IdentifierExpr{value: "a", kind: Table}},
			{source: &IdentifierExpr{value: "b", kind: Table}},
			{source: &IdentifierExpr{value: "c", kind: Table}},
		}},
		{"SELECT * FROM a JOIN b ON a.x == b.y;", []JoinedTable{
			{
				source: &IdentifierExpr{value: "a", kind: Table},
				joins: []Join{
					{
						source: JoinedTable{source: &IdentifierExpr{value: "b", kind: Table}},
						condition: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a.x", kind: Column},
							right:    &IdentifierExpr{value: "b.y", kind: Column},
						},
					},
				},
			},
		}},
		{"SELECT * FROM a LEFT JOIN b ON a.x == b.y;", []JoinedTable{
			{
				source: &IdentifierExpr{value: "a", kind: Table},
				joins: []Join{
					{
						source:   JoinedTable{source: &IdentifierExpr{value: "b", kind: Table}},
						joinType: Left,
						condition: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a.x", kind: Column},
							right:    &IdentifierExpr{value: "b.y", kind: Column},
						},
					},
				},
			},
		}},
		{"SELECT * FROM a LEFT OUTER JOIN b ON a.x == b.y;", []JoinedTable{
			{
				source: &IdentifierExpr{value: "a", kind: Table},
				joins: []Join{
					{
						source:   JoinedTable{source: &IdentifierExpr{value: "b", kind: Table}},
						joinType: LeftOuter,
						condition: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a.x", kind: Column},
							right:    &IdentifierExpr{value: "b.y", kind: Column},
						},
					},
				},
			},
		}},
		{"SELECT * FROM a JOIN b ON a.x == b.y JOIN c ON b.y == c.z;", []JoinedTable{
			{
				source: &IdentifierExpr{value: "a", kind: Table},
				joins: []Join{
					{
						source: JoinedTable{source: &IdentifierExpr{value: "b", kind: Table}},
						condition: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "a.x", kind: Column},
							right:    &IdentifierExpr{value: "b.y", kind: Column},
						},
					},
					{
						source: JoinedTable{source: &IdentifierExpr{value: "c", kind: Table}},
						condition: &BinaryExpr{
							operator: EQ,
							left:     &IdentifierExpr{value: "b.y", kind: Column},
							right:    &IdentifierExpr{value: "c.z", kind: Column},
						},
					},
				},
			},
		}},
		{"SELECT * FROM a INNER JOIN b USING(c, d);", []JoinedTable{
			{
				source: &IdentifierExpr{value: "a", kind: Table},
				joins: []Join{
					{
						source: JoinedTable{source: &IdentifierExpr{value: "b", kind: Table}},
						namedColumns: []Expr{
							&IdentifierExpr{value: "c", kind: Column},
							&IdentifierExpr{value: "d", kind: Column},
						},
					},
				},
			},
		}},
	}
	//"SELECT * FROM a LEFT JOIN (SELECT x AS y FROM b) ON c WHERE NOT(y='a');"

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if len(stmt.fromClause) != len(_case.expected) {
			t.Error("unexpected number of table expressions")
		}
		for i := 0; i < len(stmt.fromClause); i++ {
			if !eqJoinedTable(stmt.fromClause[i], _case.expected[i]) {
				t.Error("unexpected joined table")
			}
		}
	}
}

func TestParseWhereClause(t *testing.T) {
	type TestCase struct {
		statement string
		expected  Expr
	}

	var cases = [...]TestCase{
		{"SELECT * FROM test WHERE (a + b) IS NOT (b - d);", &IsExpr{
			inverse: true,
			left:    &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
			right:   &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
		},
		},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if !eqExpr(stmt.whereClause, _case.expected) {
			t.Error("unexpected expression")
		}
	}
}

func TestGroupByClause(t *testing.T) {
	type TestCase struct {
		statement string
		expected  []Expr
	}

	var cases = [...]TestCase{
		{"SELECT * FROM test GROUP BY a, b HAVING (a + b) > c;", []Expr{
			&IdentifierExpr{value: "a"},
			&IdentifierExpr{value: "b"},
		}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if len(stmt.groupByClause) != len(_case.expected) {
			t.Error("unexpected number of expressions")
		}

		for i := 0; i < len(stmt.groupByClause); i++ {
			if !eqExpr(stmt.groupByClause[i], _case.expected[i]) {
				t.Error("unexpected expression")
			}
		}
	}
}

func TestHavingClause(t *testing.T) {
	type TestCase struct {
		statement string
		expected  Expr
	}

	var cases = [...]TestCase{
		{"SELECT * FROM test GROUP BY a, b HAVING (a + b) > c;", &BinaryExpr{
			operator: GT,
			left: &BinaryExpr{
				operator: PLUS,
				left:     &IdentifierExpr{value: "a"},
				right:    &IdentifierExpr{value: "b"},
			},
			right: &IdentifierExpr{value: "c"},
		},
		},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if !eqExpr(stmt.havingClause, _case.expected) {
			t.Error("unexpected expression")
		}
	}
}

func TestOrderByClause(t *testing.T) {
	type TestCase struct {
		statement string
		expected  []OrderByExpr
	}

	var cases = [...]TestCase{
		{"SELECT * FROM test ORDER BY a ASC, (b + d) DESC;", []OrderByExpr{
			{expr: &IdentifierExpr{value: "a"}, sortOrder: ASC},
			{expr: &BinaryExpr{
				operator: PLUS,
				left:     &IdentifierExpr{value: "b"},
				right:    &IdentifierExpr{value: "d"},
			}, sortOrder: DESC},
		}},
		{"SELECT * FROM test ORDER BY (b - d) COLLATE BINARY ASC, a DESC NULLS LAST;", []OrderByExpr{
			{expr: &CollateExpr{
				collationName: "BINARY",
				expr: &BinaryExpr{
					operator: MINUS,
					left:     &IdentifierExpr{value: "b"},
					right:    &IdentifierExpr{value: "d"},
				},
			}, sortOrder: ASC},
			{expr: &IdentifierExpr{value: "a"}, sortOrder: DESC, nullsFirst: false},
		}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if len(stmt.orderByClause) != len(_case.expected) {
			t.Error("unexpected number of expressions")
		}
		for i := 0; i < len(stmt.orderByClause); i++ {
			if !eqOrderByExpr(stmt.orderByClause[i], _case.expected[i]) {
				t.Error("unexpected expression")
			}
		}
	}
}

func TestLimitClause(t *testing.T) {
	type TestCase struct {
		statement string
		expected  LimitExpr
	}

	var cases = [...]TestCase{
		{"SELECT * FROM test LIMIT (5 + 2) OFFSET (a * 6);", LimitExpr{
			count: &BinaryExpr{
				operator: PLUS,
				left:     &LiteralExpr{value: "5"},
				right:    &LiteralExpr{value: "2"},
			},
			skip: &BinaryExpr{
				operator: STAR,
				left:     &IdentifierExpr{value: "a"},
				right:    &LiteralExpr{value: "6"},
			},
		}},
		{"SELECT * FROM test LIMIT (5 + 2), (a * 6);", LimitExpr{
			count: &BinaryExpr{
				operator: PLUS,
				left:     &LiteralExpr{value: "5"},
				right:    &LiteralExpr{value: "2"},
			},
			skip: &BinaryExpr{
				operator: STAR,
				left:     &IdentifierExpr{value: "a"},
				right:    &LiteralExpr{value: "6"},
			},
		}},
	}

	for _, _case := range cases {
		stmt := parseStatement(_case.statement)
		if !eqExpr(stmt.limitClause.count, _case.expected.count) {
			t.Error("unexpected expression")
		}
		if !eqExpr(stmt.limitClause.skip, _case.expected.skip) {
			t.Error("unexpected expression")
		}
	}
}

func TestSubSelects(t *testing.T) {
	type TestCase struct {
		statement string
		expected  SelectStmt
	}

	var _ = [...]TestCase{
		/*	{"SELECT a FROM (SELECT b FROM c WHERE b == 5) AS d;", SelectStmt{
			fromClause: TableExpr{
				source: SelectStmt{
					resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "b"}}},
					fromClause:   TableExpr{source: &IdentifierExpr{value: "c"}},
					whereClause: &BinaryExpr{
						operator: EQ,
						left:     &IdentifierExpr{value: "b"},
						right:    &IdentifierExpr{value: "5"},
					},
				},
			},
		}},*/
	}

	/*	for _, testCase := range subSelectTestCases {
		stmt := parseStatement(testCase.statement)
		if _, ok := stmt.fromClause.source.(SelectStmt); !ok {
			t.Fatalf("expected SELECT statement")
		}
	}*/
}
