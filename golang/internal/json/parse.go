package json

type Parser struct {
	Tokens []Token
	pos    int
	Ast    ASTNode
}

func (p *Parser) cToken() Token {
	return p.Tokens[p.pos]
}

func (p *Parser) next() (bool, Token) {
	if p.pos >= len(p.Tokens) {
		return false, Token{}
	}

	lPos := p.pos
	p.pos++
	return true, p.Tokens[lPos]
}

func (p *Parser) parseArray(root *ASTNode) {
	for p.pos < len(p.Tokens) {
		node := ASTNode{}
		switch p.cToken().tokenType {
		case JSON_TOKEN_RIGHT_SQUARE_BRACKET:
			p.next()
			return
		case JSON_TOKEN_FALSE:
			node.Value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_NULL:
			node.Value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_TRUE:
			node.Value = JSON_VALUE_TRUE
			_, _ = p.next()
		case JSON_TOKEN_LEFT_CURLY_BRACKET:
			node.Value = JSON_VALUE_OBJECT
			p.pos++
			p.parseObject(&node)
		case JSON_TOKEN_LEFT_SQUARE_BRACKET:
			node.Value = JSON_VALUE_ARRAY
			p.pos++
			p.parseArray(&node)
		case JSON_TOKEN_NUMBER:
			node.Value = JSON_VALUE_NUMBER
			node.Number = p.cToken().number
			_, _ = p.next()
		case JSON_TOKEN_STRING:
			node.Value = JSON_VALUE_STRING
			node.String = p.cToken().string
			_, _ = p.next()
		default:
			panic("unexpected token\n")
		}

		root.Values = append(root.Values, &node)
		if p.cToken().tokenType == JSON_TOKEN_COMMA {
			p.pos++
			continue
		} else if p.cToken().tokenType == JSON_TOKEN_RIGHT_SQUARE_BRACKET {
			p.pos++
			break
		}
		panic("unexpected token")
	}
}

func (p *Parser) parseObject(root *ASTNode) {
	for p.pos < len(p.Tokens) {
		if p.cToken().tokenType == JSON_TOKEN_RIGHT_CURLY_BRACKET {
			return
		}

		member := ASTNode{}
		if p.cToken().tokenType != JSON_TOKEN_STRING {
			panic("expected a string\n")
		}
		member.Name = p.cToken().string
		p.next()

		if p.cToken().tokenType != JSON_TOKEN_COLON {
			panic("expected a separator\n")
		}
		p.next()

		switch p.cToken().tokenType {
		case JSON_TOKEN_RIGHT_CURLY_BRACKET:
			p.next()
			return
		case JSON_TOKEN_FALSE:
			member.Value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_NULL:
			member.Value = JSON_VALUE_NULL
			_, _ = p.next()
		case JSON_TOKEN_TRUE:
			member.Value = JSON_VALUE_TRUE
			_, _ = p.next()
		case JSON_TOKEN_LEFT_CURLY_BRACKET:
			member.Value = JSON_VALUE_OBJECT
			p.next()
			p.parseObject(&member)
		case JSON_TOKEN_LEFT_SQUARE_BRACKET:
			member.Value = JSON_VALUE_ARRAY
			p.next()
			p.parseArray(&member)
		case JSON_TOKEN_NUMBER:
			member.Value = JSON_VALUE_NUMBER
			member.Number = p.cToken().number
			_, _ = p.next()
		case JSON_TOKEN_STRING:
			member.Value = JSON_VALUE_STRING
			member.String = p.cToken().string
			_, _ = p.next()
		default:
			panic("unexpected token\n")
		}

		root.Members = append(root.Members, &member)
		if p.cToken().tokenType == JSON_TOKEN_COMMA {
			p.pos++
			continue
		} else if p.cToken().tokenType == JSON_TOKEN_RIGHT_CURLY_BRACKET {
			p.pos++
			break
		}
		panic("unexpected token")
	}
}

func (p *Parser) Parse() {
	switch _, token := p.next(); token.tokenType {
	case JSON_TOKEN_FALSE:
		p.Ast.Value = JSON_VALUE_FALSE
		return
	case JSON_TOKEN_NULL:
		p.Ast.Value = JSON_VALUE_NULL
		return
	case JSON_TOKEN_TRUE:
		p.Ast.Value = JSON_VALUE_TRUE
		return
	case JSON_TOKEN_LEFT_CURLY_BRACKET:
		p.Ast.Value = JSON_VALUE_OBJECT
		p.parseObject(&p.Ast)
	case JSON_TOKEN_LEFT_SQUARE_BRACKET:
		p.Ast.Value = JSON_VALUE_ARRAY
		p.parseArray(&p.Ast)
	case JSON_TOKEN_NUMBER:
		p.Ast.Value = JSON_VALUE_NUMBER
		p.Ast.Number = token.number
		return
	case JSON_TOKEN_STRING:
		p.Ast.Value = JSON_VALUE_STRING
		p.Ast.String = token.string
		return
	default:
		panic("unexpected token\n")
	}
}
