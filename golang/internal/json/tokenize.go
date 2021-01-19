package json

import (
	"strconv"
	"strings"
)

type TokenType int

const (
	JSON_TOKEN_NONE TokenType = iota
	JSON_TOKEN_LEFT_SQUARE_BRACKET
	JSON_TOKEN_LEFT_CURLY_BRACKET
	JSON_TOKEN_RIGHT_SQUARE_BRACKET
	JSON_TOKEN_RIGHT_CURLY_BRACKET
	JSON_TOKEN_COLON
	JSON_TOKEN_COMMA
	JSON_TOKEN_FALSE
	JSON_TOKEN_NULL
	JSON_TOKEN_TRUE
	JSON_TOKEN_NUMBER
	JSON_TOKEN_STRING
)

type Token struct {
	tokenType TokenType
	string    string
	number    float64
}

type Tokenizer struct {
	Buf    string
	pos    int
	Tokens []Token
}

func (t *Tokenizer) Tokenize() {
	for t.pos < len(t.Buf) {
		// Skip whitespace.
		c := t.Buf[t.pos]
		if c == 0x20 || c == 0x09 || c == 0x0A || c == 0x0D {
			t.pos++
			continue
		}

		// Handle structural characters.
		token := Token{}
		switch c {
		case '[':
			token.tokenType = JSON_TOKEN_LEFT_SQUARE_BRACKET
		case '{':
			token.tokenType = JSON_TOKEN_LEFT_CURLY_BRACKET
		case ']':
			token.tokenType = JSON_TOKEN_RIGHT_SQUARE_BRACKET
		case '}':
			token.tokenType = JSON_TOKEN_RIGHT_CURLY_BRACKET
		case ':':
			token.tokenType = JSON_TOKEN_COLON
		case ',':
			token.tokenType = JSON_TOKEN_COMMA
		}

		if token.tokenType != JSON_TOKEN_NONE {
			t.Tokens = append(t.Tokens, token)
			t.pos++
			continue
		}

		// Handle boolean literals.
		if strings.HasPrefix(t.Buf[t.pos:], "true") {
			token.tokenType = JSON_TOKEN_TRUE
			t.pos += len("true")
		} else if strings.HasPrefix(t.Buf[t.pos:], "false") {
			token.tokenType = JSON_TOKEN_FALSE
			t.pos += len("false")
		} else if strings.HasPrefix(t.Buf[t.pos:], "null") {
			token.tokenType = JSON_TOKEN_NULL
			t.pos += len("null")
		}

		if token.tokenType != JSON_TOKEN_NONE {
			t.Tokens = append(t.Tokens, token)
			continue
		}

		// Consume numeric literals.
		if (c >= '0' && c <= '9') || c == '-' {
			start := t.pos
			for {
				t.pos++
				c = t.Buf[t.pos]
				if c >= '0' && c <= '9' {
					continue
				} else if c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
					continue
				}
				break
			}
			value, _ := strconv.ParseFloat(t.Buf[start:t.pos], 64)
			token.tokenType = JSON_TOKEN_NUMBER
			token.number = value
			t.Tokens = append(t.Tokens, token)
			continue
		}

		// Must be consuming a string.
		if c == '"' {
			t.pos++
			start := t.pos
			cprev := uint8(0)
			for t.pos < len(t.Buf) {
				if t.Buf[t.pos] == '"' && cprev != '\\' {
					break
				}
				cprev = t.Buf[t.pos]
				t.pos++
			}
			token.tokenType = JSON_TOKEN_STRING
			token.string = t.Buf[start:t.pos]
			t.Tokens = append(t.Tokens, token)
			t.pos++
			continue
		}

		panic("this should never happen\n")
	}
}
