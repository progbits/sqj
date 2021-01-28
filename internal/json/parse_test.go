package json

import "testing"

type ParseTestCase struct {
	tokens []Token
	ast    ASTNode
}

var parseTestCases = []ParseTestCase{
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_STRING, string: "hello, world"},
		},
		ast: ASTNode{
			Value:  JSON_VALUE_STRING,
			String: "hello, world",
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_NUMBER, number: 3.14},
		},
		ast: ASTNode{
			Value:  JSON_VALUE_NUMBER,
			Number: 3.14,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_NULL},
		},
		ast: ASTNode{
			Value: JSON_VALUE_NULL,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_TRUE},
		},
		ast: ASTNode{
			Value: JSON_VALUE_TRUE,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_FALSE},
		},
		ast: ASTNode{
			Value: JSON_VALUE_FALSE,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_LEFT_CURLY_BRACKET},
			{tokenType: JSON_TOKEN_STRING, string: "empty"},
			{tokenType: JSON_TOKEN_COLON},
			{tokenType: JSON_TOKEN_LEFT_CURLY_BRACKET},
			{tokenType: JSON_TOKEN_RIGHT_CURLY_BRACKET},
			{tokenType: JSON_TOKEN_RIGHT_CURLY_BRACKET},
		},
		ast: ASTNode{
			Value: JSON_VALUE_OBJECT,
			Members: []*ASTNode{
				{
					Value: JSON_VALUE_OBJECT,
					Name:  "empty",
				},
			},
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_LEFT_SQUARE_BRACKET},
			{tokenType: JSON_TOKEN_LEFT_CURLY_BRACKET},
			{tokenType: JSON_TOKEN_STRING, string: "empty"},
			{tokenType: JSON_TOKEN_COLON},
			{tokenType: JSON_TOKEN_LEFT_SQUARE_BRACKET},
			{tokenType: JSON_TOKEN_RIGHT_SQUARE_BRACKET},
			{tokenType: JSON_TOKEN_RIGHT_CURLY_BRACKET},
			{tokenType: JSON_TOKEN_RIGHT_SQUARE_BRACKET},
		},
		ast: ASTNode{
			Value: JSON_VALUE_ARRAY,
			Values: []*ASTNode{
				{
					Value: JSON_VALUE_OBJECT,
					Members: []*ASTNode{
						{
							Name:  "empty",
							Value: JSON_VALUE_ARRAY,
						},
					},
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	for i := 0; i < len(parseTestCases); i++ {
		parser := Parser{
			Tokens: parseTestCases[i].tokens,
		}
		parser.Parse()

		if !equal(&parser.Ast, &parseTestCases[i].ast) {
			t.Fatal("unexpected AST")
		}
	}
}
