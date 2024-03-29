package sql

import (
	"fmt"
	"strings"
)

func precedence(token Token) int {
	switch token {
	case CONCAT:
		return 9
	case STAR, SLASH, REM:
		return 8
	case PLUS, MINUS:
		return 7
	case COLLATE:
		return 6
	case LSHIFT, RSHIFT, BITAND, BITOR:
		return 5
	case LT, LE, GT, GE:
		return 4
	case EQ, NE, IS, NOT, IN, LIKE, GLOB, MATCH, REGEXP, BETWEEN:
		return 3
	case AND:
		return 2
	case OR:
		return 1
	default:
		return 0
	}
}

// Parser is a type that converts a stream of tokens into an AST.
type Parser struct {
	scanner *Scanner
	token   Token
	value   string
}

func NewParser(scanner *Scanner) *Parser {
	return &Parser{
		scanner: scanner,
	}
}

func (p *Parser) init() {
	p.next()
}

func (p *Parser) next() {
	p.token, p.value = p.scanner.ScanToken()
}

func (p *Parser) assertAndConsumeToken(expected Token) {
	if p.token != expected {
		panic(fmt.Sprintf("unexpected token: got %s, expected %s", p.token, expected))
	}
	p.next()
}

// stmt ::= select-stmt
func (p *Parser) Parse() SelectStmt {
	p.init()

	if p.token != SELECT {
		panic("only SELECT statements are supported")
	}
	p.next()
	return p.parseSelectStmt()
}

// select-stmt ::= SELECT [ DISTINCT | ALL ]
// 			  	   result-column [, result-column ]*
//				   [ FROM table-list ]
//          	   [ WHERE expr ]
//      		   [ GROUP BY expr [, expr ]* ]
//      		   [ HAVING expr ]
//      		   [ ORDER BY ordering-term [, ordering-term ]* ]
//      		   [ LIMIT expr [ ( OFFSET | , ) expr ]
//
// TODO: Add WINDOW support.
func (p *Parser) parseSelectStmt() SelectStmt {
	stmt := SelectStmt{}
	if p.token == DISTINCT || p.token == ALL {
		stmt.isAll = p.token == ALL
		stmt.isDistinct = p.token == DISTINCT
		p.next()
	}

	// parse projection and available clauses
	stmt.resultColumn = p.parseResultColumn()
	for {
		token := p.token
		p.next()
		switch token {
		case FROM:
			stmt.fromClause = p.parseTableList()
		case WHERE:
			stmt.whereClause = p.parseExpr(0)
		case GROUP:
			p.assertAndConsumeToken(BY)
			stmt.groupByClause = append(stmt.groupByClause, p.parseExpr(0))
			for p.token == COMMA {
				p.next()
				stmt.groupByClause = append(stmt.groupByClause, p.parseExpr(0))
			}
		case HAVING:
			stmt.havingClause = p.parseExpr(0)
		case WINDOW:
			panic("WINDOW clause not currently supported")
		case ORDER:
			p.assertAndConsumeToken(BY)
			stmt.orderByClause = append(stmt.orderByClause, p.parseOrderingTerm())
			for p.token == COMMA {
				p.next()
				stmt.orderByClause = append(stmt.orderByClause, p.parseOrderingTerm())
			}
		case LIMIT:
			stmt.limitClause.count = p.parseExpr(0)
			if p.token == OFFSET || p.token == COMMA {
				p.next()
				stmt.limitClause.skip = p.parseExpr(0)
			}
		default:
			return stmt
		}
	}
}

// result-column ::= column [, column]*
func (p *Parser) parseResultColumn() []ResultColumn {
	projection := make([]ResultColumn, 0)
	projection = append(projection, p.parseColumn())
	for p.token == COMMA {
		p.next()
		projection = append(projection, p.parseColumn())
	}
	return projection
}

// column ::= * | table-name '.' '*' | expr [ AS alias ]
func (p *Parser) parseColumn() ResultColumn {
	switch p.token {
	case STAR: // *
		p.next()
		return ResultColumn{expr: &StarExpr{}}
	default: // table-name '.' '*' | expr [ AS alias ]
		expr := p.parseExpr(0)
		column := ResultColumn{expr: expr}
		if p.token == AS {
			p.next()
			column.alias = p.value
			p.next()
		}
		return column
	}
}

// table-list ::= table-or-subquery [, table-or-subquery] | join-clause
//
// join-clause ::= table_or_subquery [ join-operator table-or-subquery join-constraint ]
// join-operator ::= , | [NATURAL] [LEFT | RIGHT | FULLl] [OUTER | INNER | CROSS] JOIN
// join-args ::= [ON expr] [USING ( column-name [, column-name]* )]
func (p *Parser) parseTableList() []JoinedTable {
	tableList := make([]JoinedTable, 0)
	tableExpr := p.parseTableExpr()
	tableList = append(tableList, tableExpr)
	for {
		switch p.token {
		case COMMA:
			p.next()
			tableExpr = p.parseTableExpr()
			tableList = append(tableList, tableExpr)
			continue
		case NATURAL, INNER, LEFT, RIGHT, FULL, UNION, JOIN:
			natural := p.token == NATURAL
			var joinType JoinType
			switch p.token {
			case JOIN:
				p.next()
			case INNER:
				joinType = Inner
				p.next()
				p.next()
			case LEFT:
				p.next()
				if p.token == OUTER {
					joinType = LeftOuter
					p.next()
				} else {
					joinType = Left
				}
				p.next()
			case RIGHT:
				p.next()
				if p.token == OUTER {
					joinType = RightOuter
					p.next()
				} else {
					joinType = Right
				}
				p.next()
			case FULL:
				p.next()
				if p.token == OUTER {
					joinType = FullOuter
					p.next()
				} else {
					joinType = Full
				}
				p.next()
			}

			source := p.parseTableExpr()

			if p.token == ON {
				p.next()
				conditions := p.parseExpr(0)
				join := Join{
					source:    source,
					natural:   natural,
					joinType:  joinType,
					condition: conditions,
				}
				tableList[len(tableList)-1].joins = append(tableList[len(tableList)-1].joins, join)
				continue
			} else if p.token == USING {
				p.next()
				if p.token != LP {
					panic("expected (")
				}
				p.next()
				columns := make([]Expr, 0)
				for {
					columns = append(columns, p.parseExpr(0))
					if p.token != COMMA {
						break
					}
					p.next()
				}
				join := Join{
					source:       source,
					natural:      natural,
					joinType:     joinType,
					namedColumns: columns,
				}
				tableList[len(tableList)-1].joins = append(tableList[len(tableList)-1].joins, join)
			} else {
				panic("expected ON | USING")
			}
		default:
			return tableList
		}
		p.next()
	}
}

// table-or-subquery ::= [schema-name '.' ] table-name [AS alias]
//                       [schema-name '.' ] table-function-name ( expr [, expr]* ) [AS alias]
//						 | ( (table-or-subquery [, table-or-subquery]*) | join-clause )
//						 | (select-stmt) [AS alias]
func (p *Parser) parseTableExpr() JoinedTable {
	token, value := p.token, p.value
	switch p.next(); token {
	case IDENTIFIER:
		switch p.token {
		case AS:
			p.next()
			fallthrough
		case IDENTIFIER:
			joinedTable := JoinedTable{source: &IdentifierExpr{value: value, kind: Table, alias: p.value}}
			p.next()
			return joinedTable
		default:
			return JoinedTable{source: &IdentifierExpr{value: value, kind: Table}}
		}
	case LP:
		p.assertAndConsumeToken(SELECT)
		stmt := p.parseSelectStmt()
		// RP may not have been consumed yet.
		if p.token == RP {
			p.next()
		}

		// Consume [[AS] alias]
		if p.token == AS {
			p.next()
			p.assertAndConsumeToken(IDENTIFIER)
		} else if p.token == IDENTIFIER {
			p.next()
		}

		return JoinedTable{source: &stmt}
	default:
		panic(fmt.Sprintf("unexpected token: %s", p.token))
	}
}

// ordering-term ::= expr [COLLATE collation-name] [ ASC | DESC ]
func (p *Parser) parseOrderingTerm() OrderByExpr {
	orderingTerm := OrderByExpr{}
	orderingTerm.expr = p.parseExpr(0) // COLLATE handled here

	if p.token == ASC || p.token == DESC {
		orderingTerm.sortOrder = p.token
		p.next()
	}

	if p.token == NULLS {
		p.next()
		orderingTerm.nullsFirst = p.token == FIRST
		p.next()
	}

	return orderingTerm
}

// expr ::= literal_value
// 	    	| bind_parameter
//  	    | [ [ database_name '.' ] table_name '.' ] column_name
//		    | unary_operator expr
//		    | expr binary_operator expr
//		    | function_name '(' [ [ DISTINCT ] expr ( ',' expr ) * | '*' ] ')'
//		    | '(' expr ')'
//		    | CAST '(' expr AS type_name ')'
//		    | expr COLLATE collation_name
//		    | expr [ NOT ] ( LIKE | GLOB | REGEXP | MATCH ) expr [ ESCAPE expr ]
//		    | expr ( ISNULL | NOTNULL | NOT NULL )
//		    | expr IS [ NOT ] expr
//		    | expr [ NOT ] BETWEEN expr AND expr
//		    | expr [ NOT ] IN ( '(' [ select_stmt | expr ( ',' expr ) * ] ')' | [ database_name '.' ] table_name )
//		    | [ [ NOT ] EXISTS ] '(' select_stmt ')' - ExistsExpression
//		    | CASE [ expr ] WHEN expr THEN expr [ ELSE expr ] END
//		    | raise_function
//
// Operator precedence parsing.
func (p *Parser) parseExpr(power int) Expr {
	expr := p.parsePrefix()
	for {
		next := precedence(p.token)
		if next <= power {
			break
		}
		expr = p.parseInfix(next, expr)
	}
	return expr
}

func (p *Parser) parsePrefix() Expr {
	token, value := p.token, p.value
	switch p.next(); token {
	case IDENTIFIER:
		identifier := value
		if p.token == DOT {
			identifier += "."
			p.next()
			if p.token == IDENTIFIER || p.token == STAR {
				identifier += p.value
				p.next()
			}
		} else if p.token == LP {
			functionCallExpr := &FunctionCallExpr{function: strings.ToLower(value)}
			p.next()

			// Horrible...This needs to be a more generic method for parsing a
			// comma separated list of expressions.
			operands := p.parseResultColumn()
			for _, operand := range operands {
				functionCallExpr.operands = append(functionCallExpr.operands, operand.expr)
			}
			p.assertAndConsumeToken(RP)
			return functionCallExpr
		}
		return &IdentifierExpr{value: identifier}
	case NUMERIC_LITERAL:
		return &LiteralExpr{value: value, kind: None}
	case STRING_LITERAL:
		return &LiteralExpr{value: value, kind: Column}
	case NOT:
		if p.token == EXISTS {
			exists := p.parsePrefix()
			exists.(*ExistsExpr).inverse = true
			return exists
		}
		if p.token == LP {
			p.next()
			expr := p.parseExpr(0)
			return &UnaryExpr{operator: NOT, expr: expr}
		}
		fallthrough
	case MINUS, PLUS, BITNOT:
		return &UnaryExpr{operator: token, expr: p.parseExpr(0)}
	case CAST:
		p.assertAndConsumeToken(LP)
		left := p.parseExpr(0)
		p.assertAndConsumeToken(AS)
		alias := p.value
		p.assertAndConsumeToken(IDENTIFIER)
		p.assertAndConsumeToken(RP)
		return &CastExpr{typeName: alias, expr: left}
	case EXISTS:
		p.assertAndConsumeToken(LP)
		p.assertAndConsumeToken(SELECT)
		selectStmt := p.parseSelectStmt()
		return &ExistsExpr{selectStmt: &selectStmt}
	case CASE:
		caseExpr := CaseExpr{}
		if p.token != WHEN {
			caseExpr.expr = p.parseExpr(0)
		}

		for !(p.token == END || p.token == ELSE) {
			p.next()
			caseExpr.when = append(caseExpr.when, p.parseExpr(0))
			p.next()
			caseExpr.then = append(caseExpr.then, p.parseExpr(0))
		}

		if p.token == ELSE {
			p.next()
			caseExpr.elseExpr = p.parseExpr(0)
		}
		return &caseExpr
	case LP:
		if p.token == SELECT {
			p.next()
			stmt := p.parseSelectStmt()
			return &stmt
		}

		expr := p.parseExpr(0)
		p.assertAndConsumeToken(RP)
		return expr
	default:
		panic("unsupported expression.")
	}
}

func (p *Parser) parseInfix(power int, left Expr) Expr {
	token, _ := p.token, p.value
	switch p.next(); token {
	case COLLATE:
		collationName := p.value
		p.next() // TODO: Assert BINARY | NOCASE | RTRIM.
		return &CollateExpr{collationName: collationName, expr: left}
	case IS:
		if p.token == NOT {
			p.next()
			return &IsExpr{inverse: true, left: left, right: p.parseExpr(power)}
		}
		return &IsExpr{left: left, right: p.parseExpr(power)}
	case LIKE, GLOB, MATCH, REGEXP:
		return &StringMatchExpr{operator: token, left: left, right: p.parseExpr(power)}
	case NOT:
		switch operator := p.token; operator {
		case LIKE, GLOB, REGEXP, MATCH:
			p.next()
			right := p.parseExpr(power)
			return &StringMatchExpr{operator: operator, inverse: true, left: left, right: right}
		case BETWEEN:
			p.next()
			rangeExpr := p.parseExpr(0).(*BinaryExpr)
			return &BetweenExpr{inverse: true, expr: left, left: rangeExpr.left, right: rangeExpr.right}
		}
		right := p.parseExpr(power)
		return &BinaryExpr{operator: token, left: left, right: right}
	case BETWEEN:
		rangeExpr := p.parseExpr(0).(*BinaryExpr)
		return &BetweenExpr{expr: left, left: rangeExpr.left, right: rangeExpr.right}
	default:
		right := p.parseExpr(power)
		return &BinaryExpr{operator: token, left: left, right: right}
	}
}
