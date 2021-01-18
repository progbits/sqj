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
