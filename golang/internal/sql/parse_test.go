package sql

import (
	"testing"
)

// init a new Parser instance and parse the statement.
func parseStatement(stmt string) SelectStmt {
	scanner := NewScanner([]byte(stmt))
	parser := Parser{
		scanner: scanner,
	}
	return parser.Parse()
}

// compare two statements for equality.
func checkStatement(t *testing.T, actual, expected SelectStmt) {
	if actual.isAll != expected.isAll {
		t.Fatalf("unexpected ALL: got %v, expected %v", actual.isAll, expected.isAll)
	}
	if actual.isDistinct != expected.isDistinct {
		t.Fatalf("unexpected DISTINCT: got %v, expected %v", actual.isDistinct, expected.isDistinct)
	}

	if len(actual.resultColumn) != len(expected.resultColumn) {
		t.Fatalf("unexpected number of results: got %d, expected %d", len(actual.resultColumn), len(expected.resultColumn))
	}
	for i := 0; i < len(actual.resultColumn); i++ {
		if actual.resultColumn[i].alias != expected.resultColumn[i].alias {
			t.Fatalf("unexpected alias: got %s, expected %s", actual.resultColumn[i].alias, expected.resultColumn[i].alias)
		}
		checkExpression(t, actual.resultColumn[i].expr, expected.resultColumn[i].expr)
	}

	checkTableList(t, actual.tableList, expected.tableList)
	checkExpression(t, actual.whereClause, expected.whereClause)

	if len(actual.groupByClause) != len(expected.groupByClause) {
		t.Fatalf("unexpected number of GROUP BY expressions: got %d, expected %d", len(expected.groupByClause), len(expected.groupByClause))
	}
	for i := 0; i < len(actual.groupByClause); i++ {
		checkExpression(t, actual.groupByClause[i], expected.groupByClause[i])
	}

	checkExpression(t, actual.havingClause, expected.havingClause)

	if len(actual.orderByClause) != len(expected.orderByClause) {
		t.Fatalf("unexpected number of ORDER BY expressions: got %d, expected %d", len(actual.orderByClause), len(expected.orderByClause))
	}
	for i := 0; i < len(actual.orderByClause); i++ {
		checkOrderByExpr(t, actual.orderByClause[i], expected.orderByClause[i])
	}

	checkExpression(t, actual.limitClause.count, expected.limitClause.count)
	checkExpression(t, actual.limitClause.skip, expected.limitClause.skip)
}

// ==================
// Expression helpers
// ==================

func checkExpressions(t *testing.T, actual, expected []Expr) {
	if len(actual) != len(expected) {
		t.Fatalf("unexpected number of expressions: got %d, expected %d", len(actual), len(expected))
	}

	for i := 0; i < len(actual); i++ {
		checkExpression(t, actual[i], expected[i])
	}
}

func checkBoolean(t *testing.T, actual, expected bool) {
	if actual != expected {
		t.Fatalf("unexpected boolean value: got %v, expected %v", actual, expected)
	}
}

func checkValue(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Fatalf("unexpected value: got %v, expected %v", actual, expected)
	}
}

func checkOperator(t *testing.T, actual, expected Token) {
	if actual != expected {
		t.Fatalf("unexpected operator: got %s, expected %s", actual, expected)
	}
}

func checkExpression(t *testing.T, actual, expected Expr) {
	if actual == nil && expected == nil {
		return // base case
	}

	switch actual.(type) {
	case *StarExpr:
		if _, ok := expected.(*StarExpr); !ok {
			t.Fatalf("expected StarExpr")
		}
	case *LiteralExpr:
		checkValue(t, actual.(*LiteralExpr).value, expected.(*LiteralExpr).value)
	case *IdentifierExpr:
		checkValue(t, actual.(*IdentifierExpr).value, expected.(*IdentifierExpr).value)
	case *UnaryExpr:
		checkOperator(t, actual.(*UnaryExpr).operator, expected.(*UnaryExpr).operator)
		checkExpression(t, actual.(*UnaryExpr).expr, expected.(*UnaryExpr).expr)
	case *BinaryExpr:
		checkOperator(t, actual.(*BinaryExpr).operator, expected.(*BinaryExpr).operator)
		checkExpression(t, actual.(*BinaryExpr).left, expected.(*BinaryExpr).left)
		checkExpression(t, actual.(*BinaryExpr).right, expected.(*BinaryExpr).right)
	case *FunctionCallExpr:
		checkValue(t, actual.(*FunctionCallExpr).function, expected.(*FunctionCallExpr).function)
		checkBoolean(t, actual.(*FunctionCallExpr).distinct, expected.(*FunctionCallExpr).distinct)
		for i := 0; i < len(actual.(*FunctionCallExpr).operands); i++ {
			checkExpression(t, actual.(*FunctionCallExpr).operands[i], expected.(*FunctionCallExpr).operands[i])
		}
	case *CastExpr:
		checkValue(t, actual.(*CastExpr).typeName, expected.(*CastExpr).typeName)
		checkExpression(t, actual.(*CastExpr).expr, expected.(*CastExpr).expr)
	case *CollateExpr:
		checkValue(t, actual.(*CollateExpr).collationName, expected.(*CollateExpr).collationName)
		checkExpression(t, actual.(*CollateExpr).expr, expected.(*CollateExpr).expr)
	case *StringMatchExpr:
		checkOperator(t, actual.(*StringMatchExpr).operator, expected.(*StringMatchExpr).operator)
		checkBoolean(t, actual.(*StringMatchExpr).inverse, expected.(*StringMatchExpr).inverse)
		checkExpression(t, actual.(*StringMatchExpr).left, expected.(*StringMatchExpr).left)
		checkExpression(t, actual.(*StringMatchExpr).right, expected.(*StringMatchExpr).right)
		checkExpression(t, actual.(*StringMatchExpr).escapeExpr, expected.(*StringMatchExpr).escapeExpr)
	case *NullableExpr:
		checkOperator(t, actual.(*NullableExpr).operator, expected.(*NullableExpr).operator)
		checkExpression(t, actual.(*NullableExpr).expr, expected.(*NullableExpr).expr)
	case *IsExpr:
		checkBoolean(t, actual.(*IsExpr).inverse, expected.(*IsExpr).inverse)
		checkExpression(t, actual.(*IsExpr).left, expected.(*IsExpr).left)
		checkExpression(t, actual.(*IsExpr).right, expected.(*IsExpr).right)
	case *BetweenExpr:
		checkBoolean(t, actual.(*BetweenExpr).inverse, expected.(*BetweenExpr).inverse)
		checkExpression(t, actual.(*BetweenExpr).expr, expected.(*BetweenExpr).expr)
		checkExpression(t, actual.(*BetweenExpr).left, expected.(*BetweenExpr).left)
		checkExpression(t, actual.(*BetweenExpr).right, expected.(*BetweenExpr).right)
	case *InExpr:
		checkBoolean(t, actual.(*InExpr).inverse, expected.(*InExpr).inverse)
		checkExpression(t, actual.(*InExpr).expr, expected.(*InExpr).expr)
	case *ExistsExpr:
		checkBoolean(t, actual.(*ExistsExpr).inverse, expected.(*ExistsExpr).inverse)
		checkStatement(t, actual.(*ExistsExpr).selectStmt, expected.(*ExistsExpr).selectStmt)
	case *CaseExpr:
		checkExpression(t, actual.(*CaseExpr).expr, expected.(*CaseExpr).expr)
		if len(actual.(*CaseExpr).when) != len(expected.(*CaseExpr).when) {
			t.Fatalf("unexpected number of when expressions: got %d, expected %d", actual.(*CaseExpr).when, len(expected.(*CaseExpr).when))
		}
		checkExpressions(t, actual.(*CaseExpr).when, expected.(*CaseExpr).when)
		if len(actual.(*CaseExpr).then) != len(expected.(*CaseExpr).then) {
			t.Fatalf("unexpected number of then expressions: got %d, expected %d", actual.(*CaseExpr).then, len(expected.(*CaseExpr).then))
		}
		checkExpressions(t, actual.(*CaseExpr).then, expected.(*CaseExpr).then)
		checkExpression(t, actual.(*CaseExpr).elseExpr, expected.(*CaseExpr).elseExpr)
	case nil:
		if expected != nil {
			t.Fatalf("unexpected nil: got %v, expected %v", nil, expected)
		}
	default:
		t.Fatalf("unrecognised expression")
	}
}

// ==================
// Table list helpers
// ==================

func checkTableList(t *testing.T, actual, expected TableList) {
	if actual.source == nil && actual.source != expected.source {
		t.Fatalf("unexpected nil source")
	}

	if actual.source != nil {
		checkExpression(t, actual.source.(*IdentifierExpr), expected.source.(*IdentifierExpr))
	}

	if len(actual.joins) != len(expected.joins) {
		t.Fatalf("unexpected number of join expressions: got %d, expected %d", len(actual.joins), len(expected.joins))
	}
	for i := 0; i < len(actual.joins); i++ {
		checkJoinExpr(t, actual.joins[i], expected.joins[i])
	}
}

func checkJoinExpr(t *testing.T, actual, expected JoinExpr) {
	if actual.joinOp != expected.joinOp {
		t.Fatalf("unexpected operataor: got %v, expected %v", actual.joinOp, expected.joinOp)
	}
	checkExpression(t, actual.source.(*IdentifierExpr), expected.source.(*IdentifierExpr))
	checkJoinArgs(t, actual.joinArgs, expected.joinArgs)
}

func checkJoinArgs(t *testing.T, actual, expected JoinArgs) {
	checkExpression(t, actual.onExpr, expected.onExpr)
	if len(actual.using) != len(expected.using) {
		t.Fatalf("unexpected number of using expressions: got %d, expected %d", len(actual.using), len(expected.using))
	}
	checkExpressions(t, actual.using, expected.using)
}

// =====================
// OrderByClause Helpers
// =====================

func checkOrderByExpr(t *testing.T, actual, expected OrderByExpr) {
	checkExpression(t, actual.expr, expected.expr)
	if actual.collate != expected.collate {
		t.Fatalf("unexpected COLLATE: got %v, expected %v", actual.collate, expected.collate)
	}

	if actual.collate {
		if actual.collationName != expected.collationName {
			t.Fatalf("unexpected collation: got %v, expected %v", actual.collationName, expected.collationName)
		}
	}

	if actual.sortOrder != expected.sortOrder {
		t.Fatalf("unexpected sort order: got %v, expected %v", actual.sortOrder, expected.sortOrder)
	}

	if actual.nullsFirst != expected.nullsFirst {
		t.Fatalf("unexpected NULLS order: got %v, expected %v", actual.nullsFirst, expected.nullsFirst)
	}
}

// =====================
// Expression Test Cases
// =====================

type ExpressionTestCase struct {
	statement string
	expected  []Expr
}

var exprTestCases = [...]ExpressionTestCase{
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
	//{"SELECT a IS b;", []Expr{}{&BinaryExpr{operator: IS, left: IdentifierExpr{value: "a"}, right: IdentifierExpr{value: "b"}}}},
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

func TestParseExpr(t *testing.T) {
	for _, testCase := range exprTestCases {
		stmt := parseStatement(testCase.statement)
		if len(stmt.resultColumn) != len(testCase.expected) {
			t.Fatalf("unexpected number of columns: got %d, expected %d", len(stmt.resultColumn), len(testCase.expected))
		}
		for i := 0; i < len(stmt.resultColumn); i++ {
			checkExpression(t, stmt.resultColumn[i].expr, testCase.expected[i])
		}
	}
}

func TestParseInExpr(t *testing.T) {
	// Arrange
	stmt := "SELECT EXISTS(SELECT a FROM b WHERE a == 3);"
	expected := SelectStmt{
		resultColumn: []ResultColumn{{
			expr: &ExistsExpr{
				selectStmt: SelectStmt{
					resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "a"}}},
					tableList:    TableList{source: &IdentifierExpr{value: "b"}},
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
	checkStatement(t, parsedStmt, expected)
}

func TestParseInverseInExpr(t *testing.T) {
	// Arrange
	stmt := "SELECT NOT EXISTS(SELECT a FROM b WHERE a == 3);"
	expected := SelectStmt{
		resultColumn: []ResultColumn{{
			expr: &ExistsExpr{
				inverse: true,
				selectStmt: SelectStmt{
					resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "a"}}},
					tableList:    TableList{source: &IdentifierExpr{value: "b"}},
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
	checkStatement(t, parsedStmt, expected)
}

// =======================
// DISTINCT/ALL Test Cases
// =======================

type DistinctAllTestCase struct {
	statement string
	expected  SelectStmt
}

var distinctAllTestCases = [...]DistinctAllTestCase{
	{"SELECT DISTINCT *;", SelectStmt{isDistinct: true}},
	{"SELECT ALL *;", SelectStmt{isAll: true}},
}

func TestParseDistinctAll(t *testing.T) {
	for _, testCase := range distinctAllTestCases {
		stmt := parseStatement(testCase.statement)
		if stmt.isDistinct != testCase.expected.isDistinct {
			t.Fatalf("unexpected DISTINCT")
		}
		if stmt.isAll != testCase.expected.isAll {
			t.Fatalf("unexpected ALL")
		}
	}
}

// =========================
// Result Columns Test Cases
// =========================

type ResultColumnTestCase struct {
	statement string
	expected  []ResultColumn
}

var selectTestCases = [...]ResultColumnTestCase{
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

func TestParseResultColumns(t *testing.T) {
	for _, testCase := range selectTestCases {
		stmt := parseStatement(testCase.statement)
		if len(stmt.resultColumn) != len(testCase.expected) {
			t.Fatalf("unexpected number of result columns: got %d, expected %d", len(stmt.resultColumn), len(testCase.expected))
		}
		for i := 0; i < len(stmt.resultColumn); i++ {
			checkExpression(t, stmt.resultColumn[i].expr, testCase.expected[i].expr)
		}
	}
}

// =====================
// Table List Test Cases
// =====================

type FromTestCase struct {
	statement string
	expected  TableList
}

var tableListTestCases = [...]FromTestCase{
	{"SELECT a, b FROM a_table INNER JOIN b_table ON a_table.a == b_table.c;", TableList{
		source: &IdentifierExpr{value: "a_table"},
		joins: []JoinExpr{{
			joinOp: JoinOperator{constraint: INNER},
			source: &IdentifierExpr{value: "b_table"},
			joinArgs: JoinArgs{
				onExpr: &BinaryExpr{
					operator: EQ,
					left:     &IdentifierExpr{value: "a_table.a"},
					right:    &IdentifierExpr{value: "b_table.c"},
				},
			},
		}},
	}},
	{"SELECT * FROM a INNER JOIN b USING(c, d);", TableList{
		source: &IdentifierExpr{value: "a"},
		joins: []JoinExpr{{
			joinOp: JoinOperator{constraint: INNER},
			source: &IdentifierExpr{value: "b"},
			joinArgs: JoinArgs{
				using: []Expr{&IdentifierExpr{value: "c"}, &IdentifierExpr{value: "d"}},
			},
		}},
	}},
	{"SELECT * FROM a LEFT INNER JOIN b USING(c, d);", TableList{
		source: &IdentifierExpr{value: "a"},
		joins: []JoinExpr{{
			joinOp: JoinOperator{constraint: INNER, operator: LEFT},
			source: &IdentifierExpr{value: "b"},
			joinArgs: JoinArgs{
				using: []Expr{&IdentifierExpr{value: "c"}, &IdentifierExpr{value: "d"}},
			},
		}},
	}},
	//	{"SELECT * FROM a AS x NATURAL JOIN b;", TableList{
	//		source: &IdentifierExpr{value: "a"},
	//		joins: []JoinExpr{{
	//			joinOp: JoinOperator{constraint: NATURAL},
	//			source: &IdentifierExpr{value: "b"},
	//			joinArgs: JoinArgs{},
	//		}},
	//	}},
	//"SELECT * FROM a LEFT JOIN (SELECT x AS y FROM b) ON c WHERE NOT(y='a');"
}

func TestParseTableList(t *testing.T) {
	for _, testCase := range tableListTestCases {
		stmt := parseStatement(testCase.statement)
		checkTableList(t, stmt.tableList, testCase.expected)
	}
}

// ================
// WHERE Test Cases
// ================

type WhereTestCase struct {
	statement string
	expected  Expr
}

var whereClauseTestCases = [...]WhereTestCase{
	{"SELECT * FROM test WHERE (a + b) IS NOT (b - d);", &IsExpr{
		inverse: true,
		left:    &BinaryExpr{operator: PLUS, left: &IdentifierExpr{value: "a"}, right: &IdentifierExpr{value: "b"}},
		right:   &BinaryExpr{operator: MINUS, left: &IdentifierExpr{value: "b"}, right: &IdentifierExpr{value: "d"}},
	},
	},
}

func TestParseWhereClause(t *testing.T) {
	for _, testCase := range whereClauseTestCases {
		stmt := parseStatement(testCase.statement)
		checkExpression(t, stmt.whereClause, testCase.expected)
	}
}

// ===================
// GROUP BY Test Cases
// ===================

type GroupByTestCase struct {
	statement string
	expected  []Expr
}

var groupByClauseTestCases = [...]GroupByTestCase{
	{"SELECT * FROM test GROUP BY a, b HAVING (a + b) > c;", []Expr{
		&IdentifierExpr{value: "a"},
		&IdentifierExpr{value: "b"},
	}},
}

func TestGroupByClause(t *testing.T) {
	for _, testCase := range groupByClauseTestCases {
		stmt := parseStatement(testCase.statement)
		checkExpressions(t, stmt.groupByClause, testCase.expected)
	}
}

// =================
// HAVING Test Cases
// =================

type HavingTestCase struct {
	statement string
	expected  Expr
}

var havingClauseTestCases = [...]HavingTestCase{
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

func TestHavingClause(t *testing.T) {
	for _, testCase := range havingClauseTestCases {
		stmt := parseStatement(testCase.statement)
		checkExpression(t, stmt.havingClause, testCase.expected)
	}
}

// ===================
// ORDER BY Test Cases
// ===================

type OrderByTestCase struct {
	statement string
	expected  []OrderByExpr
}

var orderByClauseTestCases = [...]OrderByTestCase{
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

func TestOrderByClause(t *testing.T) {
	for _, testCase := range orderByClauseTestCases {
		stmt := parseStatement(testCase.statement)
		if len(stmt.orderByClause) != len(testCase.expected) {
			t.Fatalf("unexpected number of expressions: got %d, expected %d", len(stmt.orderByClause), len(testCase.expected))
		}
		for i := 0; i < len(stmt.orderByClause); i++ {
			checkOrderByExpr(t, stmt.orderByClause[i], testCase.expected[i])
		}
	}
}

// ================
// LIMIT Test Cases
// ================

type LimitTestCase struct {
	statement string
	expected  LimitExpr
}

var limitClauseTestCases = [...]LimitTestCase{
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

func TestLimitClause(t *testing.T) {
	for _, testCase := range limitClauseTestCases {
		stmt := parseStatement(testCase.statement)
		checkExpression(t, stmt.limitClause.count, testCase.expected.count)
		checkExpression(t, stmt.limitClause.skip, testCase.expected.skip)
	}
}

// ===========
// Sub-selects
// ===========

type SubSelectTestCase struct {
	statement string
	expected  SelectStmt
}

var subSelectTestCases = [...]SubSelectTestCase{
	{"SELECT a FROM (SELECT b FROM c WHERE b == 5) AS d;", SelectStmt{
		tableList: TableList{
			source: SelectStmt{
				resultColumn: []ResultColumn{{expr: &IdentifierExpr{value: "b"}}},
				tableList:    TableList{source: &IdentifierExpr{value: "c"}},
				whereClause: &BinaryExpr{
					operator: EQ,
					left:     &IdentifierExpr{value: "b"},
					right:    &IdentifierExpr{value: "5"},
				},
			},
		},
	}},
}

func TestSubSelects(t *testing.T) {
	for _, testCase := range subSelectTestCases {
		stmt := parseStatement(testCase.statement)
		if _, ok := stmt.tableList.source.(SelectStmt); !ok {
			t.Fatalf("expected SELECT statement")
		}
	}
}
