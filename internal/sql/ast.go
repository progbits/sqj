package sql

import "sort"

type IdentifierKind int

const (
	Column IdentifierKind = iota
	Table
)

type JoinType int

const (
	Inner JoinType = iota
	Left
	LeftOuter
	Right
	RightOuter
	Full
	FullOuter
	Union
)

// Empty table expression interface.
type TableExpr interface {
	isTableExpr()
}

func (t *SelectStmt) isTableExpr()     {}
func (t *IdentifierExpr) isTableExpr() {}

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
func (e *SelectStmt) isExpr()       {}

// expression types
type (
	StarExpr struct{}

	LiteralExpr struct {
		value string
	}

	IdentifierExpr struct {
		value string
		alias string
		kind  IdentifierKind
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
		selectStmt *SelectStmt
	}

	CaseExpr struct {
		expr     Expr
		when     []Expr
		then     []Expr
		elseExpr Expr
	}
)

// eqExpr checks two expressions for equality.
func eqExpr(a, b Expr) bool {
	if a == nil && b == nil {
		return true
	}

	switch a.(type) {
	case *StarExpr:
		if _, ok := b.(*StarExpr); !ok {
			return false
		}
		return true
	case *LiteralExpr:
		if _, ok := b.(*LiteralExpr); !ok {
			return false
		}
		return a.(*LiteralExpr).value == b.(*LiteralExpr).value
	case *IdentifierExpr:
		if _, ok := b.(*IdentifierExpr); !ok {
			return false
		}

		if a.(*IdentifierExpr).value != b.(*IdentifierExpr).value {
			return false
		}
		return a.(*IdentifierExpr).kind == b.(*IdentifierExpr).kind
	case *UnaryExpr:
		if _, ok := b.(*UnaryExpr); !ok {
			return false
		}

		if a.(*UnaryExpr).operator != b.(*UnaryExpr).operator {
			return false
		}
		return eqExpr(a.(*UnaryExpr).expr, b.(*UnaryExpr).expr)
	case *BinaryExpr:
		if _, ok := b.(*BinaryExpr); !ok {
			return false
		}

		if a.(*BinaryExpr).operator != b.(*BinaryExpr).operator {
			return false
		}

		if !eqExpr(a.(*BinaryExpr).left, b.(*BinaryExpr).left) {
			return false
		}
		return eqExpr(a.(*BinaryExpr).right, b.(*BinaryExpr).right)
	case *FunctionCallExpr:
		if _, ok := b.(*FunctionCallExpr); !ok {
			return false
		}

		if a.(*FunctionCallExpr).function != b.(*FunctionCallExpr).function {
			return false
		}

		if a.(*FunctionCallExpr).distinct != b.(*FunctionCallExpr).distinct {
			return false
		}

		if len(a.(*FunctionCallExpr).operands) != len(b.(*FunctionCallExpr).operands) {
			return false
		}
		for i := 0; i < len(a.(*FunctionCallExpr).operands); i++ {
			if !eqExpr(a.(*FunctionCallExpr).operands[i], b.(*FunctionCallExpr).operands[i]) {
				return false
			}
		}
		return true
	case *CastExpr:
		if _, ok := b.(*CastExpr); !ok {
			return false
		}

		if a.(*CastExpr).typeName != b.(*CastExpr).typeName {
			return false
		}
		return eqExpr(a.(*CastExpr).expr, b.(*CastExpr).expr)
	case *CollateExpr:
		if _, ok := b.(*CollateExpr); !ok {
			return false
		}

		if a.(*CollateExpr).collationName != b.(*CollateExpr).collationName {
			return false
		}
		return eqExpr(a.(*CollateExpr).expr, b.(*CollateExpr).expr)
	case *StringMatchExpr:
		if _, ok := b.(*StringMatchExpr); !ok {
			return false
		}

		if a.(*StringMatchExpr).operator != b.(*StringMatchExpr).operator {
			return false
		}

		if a.(*StringMatchExpr).inverse != b.(*StringMatchExpr).inverse {
			return false
		}

		if !eqExpr(a.(*StringMatchExpr).left, b.(*StringMatchExpr).left) {
			return false
		}

		if !eqExpr(a.(*StringMatchExpr).right, b.(*StringMatchExpr).right) {
			return false
		}
		return eqExpr(a.(*StringMatchExpr).escapeExpr, b.(*StringMatchExpr).escapeExpr)
	case *NullableExpr:
		if _, ok := b.(*NullableExpr); !ok {
			return false
		}

		if a.(*NullableExpr).operator != b.(*NullableExpr).operator {
			return false
		}
		return eqExpr(a.(*NullableExpr).expr, b.(*NullableExpr).expr)
	case *IsExpr:
		if _, ok := b.(*IsExpr); !ok {
			return false
		}

		if a.(*IsExpr).inverse != b.(*IsExpr).inverse {
			return false
		}

		if !eqExpr(a.(*IsExpr).left, b.(*IsExpr).left) {
			return false
		}

		return eqExpr(a.(*IsExpr).right, b.(*IsExpr).right)
	case *BetweenExpr:
		if _, ok := b.(*BetweenExpr); !ok {
			return false
		}

		if a.(*BetweenExpr).inverse != b.(*BetweenExpr).inverse {
			return false
		}

		if !eqExpr(a.(*BetweenExpr).expr, b.(*BetweenExpr).expr) {
			return false
		}

		if !eqExpr(a.(*BetweenExpr).left, b.(*BetweenExpr).left) {
			return false
		}
		return eqExpr(a.(*BetweenExpr).right, b.(*BetweenExpr).right)
	case *InExpr:
		if _, ok := b.(*InExpr); !ok {
			return false
		}

		if a.(*InExpr).inverse != b.(*InExpr).inverse {
			return false
		}
		return eqExpr(a.(*InExpr).expr, b.(*InExpr).expr)
	case *ExistsExpr:
		if _, ok := b.(*ExistsExpr); !ok {
			return false
		}

		if a.(*ExistsExpr).inverse != b.(*ExistsExpr).inverse {
			return false
		}
		return eqSelectStmt(a.(*ExistsExpr).selectStmt, b.(*ExistsExpr).selectStmt)
	case *CaseExpr:
		if _, ok := b.(*CaseExpr); !ok {
			return false
		}

		if !eqExpr(a.(*CaseExpr).expr, b.(*CaseExpr).expr) {
			return false
		}

		if len(a.(*CaseExpr).when) != len(b.(*CaseExpr).when) {
			return false
		}
		for i := 0; i < len(a.(*CaseExpr).when); i++ {
			if !eqExpr(a.(*CaseExpr).when[i], b.(*CaseExpr).when[i]) {
				return false
			}
		}

		if len(a.(*CaseExpr).then) != len(b.(*CaseExpr).then) {
			return false
		}
		for i := 0; i < len(a.(*CaseExpr).then); i++ {
			if !eqExpr(a.(*CaseExpr).then[i], b.(*CaseExpr).then[i]) {
				return false
			}
		}
		return eqExpr(a.(*CaseExpr).elseExpr, b.(*CaseExpr).elseExpr)
	case nil:
		if b != nil {
			return false
		}
	}
	return false
}

type ResultColumn struct {
	alias string
	expr  Expr
}

// eqResultColumn checks two ResultColumns for equality
func eqResultColumn(a, b ResultColumn) bool {
	if a.alias != b.alias {
		return false
	}
	return eqExpr(a.expr, b.expr)
}

type JoinedTable struct {
	source TableExpr
	joins  []Join
}

// eqTableExpr checks two TableExprs for equality.
func eqTableExpr(a, b TableExpr) bool {
	if a == nil && b == nil {
		return true
	}

	switch a.(type) {
	case *SelectStmt:
		if _, ok := b.(*SelectStmt); !ok {
			return false
		}
		return eqSelectStmt(a.(*SelectStmt), b.(*SelectStmt))
	case *IdentifierExpr:
		if _, ok := b.(*IdentifierExpr); !ok {
			return false
		}
		return eqExpr(a.(*IdentifierExpr), b.(*IdentifierExpr))

	default:
		return false
	}
}

// eqJoinedTable checks two JoinedTables for equality.
func eqJoinedTable(a, b JoinedTable) bool {
	if !eqTableExpr(a.source, b.source) {
		return false
	}

	if len(a.joins) != len(b.joins) {
		return false
	}

	for i := 0; i < len(a.joins); i++ {
		if !eqJoin(a.joins[i], b.joins[i]) {
			return false
		}
	}
	return true
}

type Join struct {
	source       JoinedTable
	natural      bool
	joinType     JoinType
	condition    Expr
	namedColumns []Expr
}

// eqJoin checks two Joins for equality.
func eqJoin(a, b Join) bool {
	if !eqJoinedTable(a.source, b.source) {
		return false
	}

	if a.natural != b.natural {
		return false
	}

	if a.joinType != b.joinType {
		return false
	}

	if !eqExpr(a.condition, b.condition) {
		return false
	}

	if len(a.namedColumns) != len(b.namedColumns) {
		return false
	}
	for i := 0; i < len(a.namedColumns); i++ {
		if !eqExpr(a.namedColumns[i], b.namedColumns[i]) {
			return false
		}
	}

	return true
}

type OrderByExpr struct {
	expr          Expr
	collate       bool
	collationName string
	sortOrder     Token // ASC | DESC
	nullsFirst    bool
}

// eqOrderByExpr checks two OrderByExpr for equality.
func eqOrderByExpr(a, b OrderByExpr) bool {
	if !eqExpr(a.expr, b.expr) {
		return false
	}

	if a.collate != b.collate {
		return false
	}

	if a.collationName != b.collationName {
		return false
	}

	if a.sortOrder != b.sortOrder {
		return false
	}

	if a.nullsFirst != b.nullsFirst {
		return false
	}

	return true
}

type LimitExpr struct {
	count Expr
	skip  Expr
}

// eqLimitExpr checks two OrderByExpr for equality.
func eqLimitExpr(a, b LimitExpr) bool {
	if !eqExpr(a.count, b.count) {
		return false
	}
	return eqExpr(a.skip, b.skip)
}

type SelectStmt struct {
	isAll      bool
	isDistinct bool

	resultColumn  []ResultColumn
	fromClause    []JoinedTable
	whereClause   Expr
	groupByClause []Expr
	havingClause  Expr
	orderByClause []OrderByExpr
	limitClause   LimitExpr
}

// eqSelectStmt checks two SelectStmts for equality.
func eqSelectStmt(a, b *SelectStmt) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil && a != b {
		return false
	}

	if a.isAll != b.isAll {
		return false
	}

	if a.isDistinct != b.isDistinct {
		return false
	}

	if len(a.resultColumn) != len(b.resultColumn) {
		return false
	}
	for i := 0; i < len(a.resultColumn); i++ {
		if !eqResultColumn(a.resultColumn[i], b.resultColumn[i]) {
			return false
		}
	}

	if len(a.fromClause) != len(b.fromClause) {
		return false
	}
	for i := 0; i < len(a.fromClause); i++ {
		if !eqJoinedTable(a.fromClause[i], b.fromClause[i]) {
			return false
		}
	}

	if !eqExpr(a.whereClause, b.whereClause) {
		return false
	}

	if len(a.groupByClause) != len(b.groupByClause) {
		return false
	}
	for i := 0; i < len(a.groupByClause); i++ {
		if !eqExpr(a.groupByClause[i], b.groupByClause[i]) {
			return false
		}
	}

	if !eqExpr(a.havingClause, b.havingClause) {
		return false
	}

	if len(a.orderByClause) != len(b.orderByClause) {
		return false
	}
	for i := 0; i < len(a.orderByClause); i++ {
		if !eqOrderByExpr(a.orderByClause[i], b.orderByClause[i]) {
			return false
		}
	}

	return eqLimitExpr(a.limitClause, b.limitClause)
}

// extractIdentifierFromExpression returns all identifiers present in an expression
func extractIdentifierFromExpression(expr Expr, kind IdentifierKind, idents map[string]bool) {
	switch expr.(type) {
	case *StarExpr:
		if kind != Table {
			if found := idents["*"]; !found {
				idents["*"] = true
			}
		}
	case *LiteralExpr:
		break
	case *IdentifierExpr:
		value := expr.(*IdentifierExpr)
		if value.kind == kind {
			if found := idents[expr.(*IdentifierExpr).value]; !found {
				idents[expr.(*IdentifierExpr).value] = true
			}
		}
	case *UnaryExpr:
		extractIdentifierFromExpression(expr.(*UnaryExpr).expr, kind, idents)
	case *BinaryExpr:
		extractIdentifierFromExpression(expr.(*BinaryExpr).left, kind, idents)
		extractIdentifierFromExpression(expr.(*BinaryExpr).right, kind, idents)
	case *FunctionCallExpr:
		funCallExpr := expr.(*FunctionCallExpr)
		for i := 0; i < len(funCallExpr.operands); i++ {
			extractIdentifierFromExpression(funCallExpr.operands[i], kind, idents)
		}
	case *CastExpr:
		extractIdentifierFromExpression(expr.(*CastExpr).expr, kind, idents)
	case *CollateExpr:
		extractIdentifierFromExpression(expr.(*CollateExpr).expr, kind, idents)
	case *StringMatchExpr:
		extractIdentifierFromExpression(expr.(*StringMatchExpr).left, kind, idents)
		extractIdentifierFromExpression(expr.(*StringMatchExpr).right, kind, idents)
		extractIdentifierFromExpression(expr.(*StringMatchExpr).escapeExpr, kind, idents)
	case *NullableExpr:
		extractIdentifierFromExpression(expr.(*NullableExpr).expr, kind, idents)
	case *IsExpr:
		extractIdentifierFromExpression(expr.(*IsExpr).left, kind, idents)
		extractIdentifierFromExpression(expr.(*IsExpr).right, kind, idents)
	case *BetweenExpr:
		extractIdentifierFromExpression(expr.(*BetweenExpr).expr, kind, idents)
		extractIdentifierFromExpression(expr.(*BetweenExpr).left, kind, idents)
		extractIdentifierFromExpression(expr.(*BetweenExpr).right, kind, idents)
	case *InExpr:
		extractIdentifierFromExpression(expr.(*InExpr).expr, kind, idents)
	case *ExistsExpr:
		extractIdentifiersImpl(expr.(*ExistsExpr).selectStmt, kind, idents)
	case *CaseExpr:
		caseExpr := expr.(*CaseExpr)
		extractIdentifierFromExpression(caseExpr.expr, kind, idents)
		for i := 0; i < len(caseExpr.when); i++ {
			extractIdentifierFromExpression(caseExpr.when[i], kind, idents)
		}
		for i := 0; i < len(caseExpr.then); i++ {
			extractIdentifierFromExpression(caseExpr.then[i], kind, idents)
		}
		extractIdentifierFromExpression(caseExpr.elseExpr, kind, idents)
	case *SelectStmt:
		selectStmt := expr.(*SelectStmt)
		extractIdentifiersImpl(selectStmt, kind, idents)
	}
}

func extractIndentifiersFromJoinedTable(joinedTable *JoinedTable, kind IdentifierKind, idents map[string]bool) {
	switch joinedTable.source.(type) {
	case Expr:
		extractIdentifierFromExpression(joinedTable.source.(Expr), kind, idents)
	case *SelectStmt:
		selectStmt := joinedTable.source.(*SelectStmt)
		extractIdentifiersImpl(selectStmt, kind, idents)
	default:
		panic("unexpected table list source")
	}

	for i := 0; i < len(joinedTable.joins); i++ {
		extractIndentifiersFromJoinedTable(&joinedTable.joins[i].source, kind, idents)
		extractIdentifierFromExpression(joinedTable.joins[i].condition, kind, idents)
		for k := 0; k < len(joinedTable.joins[i].namedColumns); k++ {
			columnExpr := joinedTable.joins[i].namedColumns[k]
			extractIdentifierFromExpression(columnExpr, kind, idents)
		}
	}
}

func extractIdentifiersImpl(stmt *SelectStmt, kind IdentifierKind, idents map[string]bool) {
	// Extract identifiers from result columns.
	for i := 0; i < len(stmt.resultColumn); i++ {
		extractIdentifierFromExpression(stmt.resultColumn[i].expr, kind, idents)
	}

	// Handle the case where our table list either an identifier or a sub-query.
	for i := 0; i < len(stmt.fromClause); i++ {
		extractIndentifiersFromJoinedTable(&stmt.fromClause[i], kind, idents)
	}

	// Handle WHERE clause expressions.
	extractIdentifierFromExpression(stmt.whereClause, kind, idents)

	// Handle GROUP BY clause.
	for i := 0; i < len(stmt.groupByClause); i++ {
		extractIdentifierFromExpression(stmt.groupByClause[i], kind, idents)
	}

	// Handle HAVING clause.
	extractIdentifierFromExpression(stmt.havingClause, kind, idents)

	// Handle ORDER BY clauses.
	for i := 0; i < len(stmt.orderByClause); i++ {
		extractIdentifierFromExpression(stmt.orderByClause[i].expr, kind, idents)
	}

	// Handle LIMIT clause.
	extractIdentifierFromExpression(stmt.limitClause.skip, kind, idents)
	extractIdentifierFromExpression(stmt.limitClause.count, kind, idents)
}

// ExtractIdentifiers returns all identifiers from a SELECT statement.
func ExtractIdentifiers(stmt *SelectStmt, kind IdentifierKind) []string {
	uniqueIdentifiers := make(map[string]bool, 0)
	extractIdentifiersImpl(stmt, kind, uniqueIdentifiers)

	identifiers := make([]string, 0)
	for key, _ := range uniqueIdentifiers {
		identifiers = append(identifiers, key)
	}

	sort.Strings(identifiers)
	return identifiers
}
