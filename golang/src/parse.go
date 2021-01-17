package main

type Parser struct {
	tokens []Token
	pos    int
	ast    ASTNode
}

func (p *Parser) cToken() Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() (bool, Token) {
	if p.pos >= len(p.tokens) {
		return false, Token{}
	}

	lPos := p.pos
	p.pos++
	return true, p.tokens[lPos]
}

func (p *Parser) parseArray(root *ASTNode) {
	for p.pos < len(p.tokens) {
		node := ASTNode{}
		switch p.cToken().tokenType {
		case JSON_TOKEN_RIGHT_SQUARE_BRACKET:
			p.next()
			return
		case JSON_TOKEN_FALSE:
			node.value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_NULL:
			node.value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_TRUE:
			node.value = JSON_VALUE_TRUE
			_, _ = p.next()
		case JSON_TOKEN_LEFT_CURLY_BRACKET:
			node.value = JSON_VALUE_OBJECT
			p.pos++
			p.parseObject(&node)
		case JSON_TOKEN_LEFT_SQUARE_BRACKET:
			node.value = JSON_VALUE_ARRAY
			p.pos++
			p.parseArray(&node)
		case JSON_TOKEN_NUMBER:
			node.value = JSON_VALUE_NUMBER
			node.number = p.cToken().number
			_, _ = p.next()
		case JSON_TOKEN_STRING:
			node.value = JSON_VALUE_STRING
			node.string = p.cToken().string
			_, _ = p.next()
		default:
			panic("unexpected token\n")
		}

		root.values = append(root.values, &node)
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
	for p.pos < len(p.tokens) {
		if p.cToken().tokenType == JSON_TOKEN_RIGHT_CURLY_BRACKET {
			return
		}

		member := ASTNode{}
		if p.cToken().tokenType != JSON_TOKEN_STRING {
			panic("expected a string\n")
		}
		member.name = p.cToken().string
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
			member.value = JSON_VALUE_FALSE
			_, _ = p.next()
		case JSON_TOKEN_NULL:
			member.value = JSON_VALUE_NULL
			_, _ = p.next()
		case JSON_TOKEN_TRUE:
			member.value = JSON_VALUE_TRUE
			_, _ = p.next()
		case JSON_TOKEN_LEFT_CURLY_BRACKET:
			member.value = JSON_VALUE_OBJECT
			p.next()
			p.parseObject(&member)
		case JSON_TOKEN_LEFT_SQUARE_BRACKET:
			member.value = JSON_VALUE_ARRAY
			p.next()
			p.parseArray(&member)
		case JSON_TOKEN_NUMBER:
			member.value = JSON_VALUE_NUMBER
			member.number = p.cToken().number
			_, _ = p.next()
		case JSON_TOKEN_STRING:
			member.value = JSON_VALUE_STRING
			member.string = p.cToken().string
			_, _ = p.next()
		default:
			panic("unexpected token\n")
		}

		root.members = append(root.members, &member)
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

func (p *Parser) parse() {
	switch _, token := p.next(); token.tokenType {
	case JSON_TOKEN_FALSE:
		p.ast.value = JSON_VALUE_FALSE
		return
	case JSON_TOKEN_NULL:
		p.ast.value = JSON_VALUE_NULL
		return
	case JSON_TOKEN_TRUE:
		p.ast.value = JSON_VALUE_TRUE
		return
	case JSON_TOKEN_LEFT_CURLY_BRACKET:
		p.ast.value = JSON_VALUE_OBJECT
		p.parseObject(&p.ast)
	case JSON_TOKEN_LEFT_SQUARE_BRACKET:
		p.ast.value = JSON_VALUE_ARRAY
		p.parseArray(&p.ast)
	case JSON_TOKEN_NUMBER:
		p.ast.value = JSON_VALUE_NUMBER
		p.ast.number = token.number
		return
	case JSON_TOKEN_STRING:
		p.ast.value = JSON_VALUE_STRING
		p.ast.string = token.string
		return
	default:
		panic("unexpected token\n")
	}
}
