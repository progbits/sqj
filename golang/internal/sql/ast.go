package sql

// Empty expression interface.
type Expr interface {
	isExpr()
}

func (e *StarExpr) isExpr()         {}
func (e *LiteralExpr) isExpr()      {}
func (e *IdentifierExpr) isExpr()   {}
func (e *UnaryExpr) isExpr()        {}
func (e *BinaryExpr) isExpr()       {}
func (e *FunctionCallExpr) isExpr() {}
func (e *CastExpr) isExpr()         {}
func (e *CollateExpr) isExpr()      {}
func (e *StringMatchExpr) isExpr()  {}
func (e *NullableExpr) isExpr()     {}
func (e *IsExpr) isExpr()           {}
func (e *BetweenExpr) isExpr()      {}
func (e *InExpr) isExpr()           {}
func (e *ExistsExpr) isExpr()       {}
func (e *CaseExpr) isExpr()         {}

// expression types
type (
	StarExpr struct{}

	LiteralExpr struct {
		value string
	}

	IdentifierExpr struct {
		value string
	}

	UnaryExpr struct {
		operator Token
		expr     Expr
	}

	BinaryExpr struct {
		operator Token
		left     Expr
		right    Expr
	}

	FunctionCallExpr struct {
		function string
		distinct bool
		operands []Expr
	}

	CastExpr struct {
		typeName string
		expr     Expr
	}

	CollateExpr struct {
		collationName string
		expr          Expr
	}

	StringMatchExpr struct {
		operator   Token
		inverse    bool
		left       Expr
		right      Expr
		escapeExpr Expr
	}

	NullableExpr struct {
		operator Token
		expr     Expr
	}

	IsExpr struct {
		inverse bool
		left    Expr
		right   Expr
	}

	BetweenExpr struct {
		inverse bool
		expr    Expr
		left    Expr
		right   Expr
	}

	InExpr struct {
		inverse bool
		expr    Expr
	}

	ExistsExpr struct {
		inverse    bool
		selectStmt SelectStmt
	}

	CaseExpr struct {
		expr     Expr
		when     []Expr
		then     []Expr
		elseExpr Expr
	}
)

type ResultColumn struct {
	alias string
	expr  Expr
}

type JoinOperator struct {
	natural    bool
	operator   Token
	constraint Token
}

type JoinArgs struct {
	onExpr Expr
	using  []Expr // IdentifierExpr
}

type JoinExpr struct {
	joinOp   JoinOperator
	source   interface{} // table or sub-select
	joinArgs JoinArgs
}

type TableList struct {
	source interface{} // table or sub-select
	joins  []JoinExpr
}

type OrderByExpr struct {
	expr          Expr
	collate       bool
	collationName string
	sortOrder     Token // ASC | DESC
	nullsFirst    bool
}

type LimitExpr struct {
	count Expr
	skip  Expr
}

type SelectStmt struct {
	isAll      bool
	isDistinct bool

	resultColumn  []ResultColumn
	tableList     TableList
	whereClause   Expr
	groupByClause []Expr
	havingClause  Expr
	orderByClause []OrderByExpr
	limitClause   LimitExpr
}

// extractIdentifierFromExpression returns all identifiers present in an expression
func extractIdentifierFromExpression(expr Expr, idents *[]string) {
	switch expr.(type) {
	case *StarExpr:
		*idents = append(*idents, "*")
	case *LiteralExpr:
		*idents = append(*idents, expr.(*LiteralExpr).value)
	case *IdentifierExpr:
		*idents = append(*idents, expr.(*IdentifierExpr).value)
	case *UnaryExpr:
		extractIdentifierFromExpression(expr.(*UnaryExpr).expr, idents)
	case *BinaryExpr:
		extractIdentifierFromExpression(expr.(*BinaryExpr).left, idents)
		extractIdentifierFromExpression(expr.(*BinaryExpr).right, idents)
	case *FunctionCallExpr:
		funCallExpr := expr.(*FunctionCallExpr)
		for i := 0; i < len(funCallExpr.operands); i++ {
			extractIdentifierFromExpression(funCallExpr.operands[i], idents)
		}
	case *CastExpr:
		extractIdentifierFromExpression(expr.(*CastExpr).expr, idents)
	case *CollateExpr:
		extractIdentifierFromExpression(expr.(*CollateExpr).expr, idents)
	case *StringMatchExpr:
		extractIdentifierFromExpression(expr.(*StringMatchExpr).left, idents)
		extractIdentifierFromExpression(expr.(*StringMatchExpr).right, idents)
		extractIdentifierFromExpression(expr.(*StringMatchExpr).escapeExpr, idents)
	case *NullableExpr:
		extractIdentifierFromExpression(expr.(*NullableExpr).expr, idents)
	case *IsExpr:
		extractIdentifierFromExpression(expr.(*IsExpr).left, idents)
		extractIdentifierFromExpression(expr.(*IsExpr).right, idents)
	case *BetweenExpr:
		extractIdentifierFromExpression(expr.(*BetweenExpr).expr, idents)
		extractIdentifierFromExpression(expr.(*BetweenExpr).left, idents)
		extractIdentifierFromExpression(expr.(*BetweenExpr).right, idents)
	case *InExpr:
		extractIdentifierFromExpression(expr.(*InExpr).expr, idents)
	case *ExistsExpr:
		extractIdentifiersImpl(&expr.(*ExistsExpr).selectStmt, idents)
	case *CaseExpr:
		caseExpr := expr.(*CaseExpr)
		extractIdentifierFromExpression(caseExpr.expr, idents)
		for i := 0; i < len(caseExpr.when); i++ {
			extractIdentifierFromExpression(caseExpr.when[i], idents)
		}
		for i := 0; i < len(caseExpr.then); i++ {
			extractIdentifierFromExpression(caseExpr.then[i], idents)
		}
		extractIdentifierFromExpression(caseExpr.elseExpr, idents)
	}
}

func extractIdentifiersImpl(stmt *SelectStmt, idents *[]string) {
	// Extract identifiers from result columns.
	for i := 0; i < len(stmt.resultColumn); i++ {
		extractIdentifierFromExpression(stmt.resultColumn[i].expr, idents)
	}

	// Handle the case where our table list either an identifier or a sub-query.
	// TODO: Handle TableList joins.
	switch stmt.tableList.source.(type) {
	case Expr:
		extractIdentifierFromExpression(stmt.tableList.source.(Expr), idents)
	case SelectStmt:
		extractIdentifiersImpl(stmt.tableList.source.(*SelectStmt), idents)
	default:
		panic("unexpected table list source")
	}

	// Handle WHERE clause expressions.
	extractIdentifierFromExpression(stmt.whereClause, idents)

	// Handle GROUP BY clause.
	for i := 0; i < len(stmt.groupByClause); i++ {
		extractIdentifierFromExpression(stmt.groupByClause[i], idents)
	}

	// Handle HAVING clause.
	extractIdentifierFromExpression(stmt.havingClause, idents)

	// Handle ORDER BY clauses.
	for i := 0; i < len(stmt.orderByClause); i++ {
		extractIdentifierFromExpression(stmt.orderByClause[i].expr, idents)
	}

	// Handle LIMIT clause.
	extractIdentifierFromExpression(stmt.limitClause.skip, idents)
	extractIdentifierFromExpression(stmt.limitClause.count, idents)
}

// ExtractIdentifiers returns all identifiers from a SELECT statement.
func ExtractIdentifiers(stmt *SelectStmt) []string {
	idents := make([]string, 0)
	extractIdentifiersImpl(stmt, &idents)

	return idents
}
