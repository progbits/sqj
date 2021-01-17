package main

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
			value:  JSON_VALUE_STRING,
			string: "hello, world",
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_NUMBER, number: 3.14},
		},
		ast: ASTNode{
			value:  JSON_VALUE_NUMBER,
			number: 3.14,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_NULL},
		},
		ast: ASTNode{
			value: JSON_VALUE_NULL,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_TRUE},
		},
		ast: ASTNode{
			value: JSON_VALUE_TRUE,
		},
	},
	{
		tokens: []Token{
			{tokenType: JSON_TOKEN_FALSE},
		},
		ast: ASTNode{
			value: JSON_VALUE_FALSE,
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
			value: JSON_VALUE_OBJECT,
			members: []*ASTNode{
				{
					value: JSON_VALUE_OBJECT,
					name:  "empty",
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
			value: JSON_VALUE_ARRAY,
			values: []*ASTNode{
				{
					value: JSON_VALUE_OBJECT,
					members: []*ASTNode{
						{
							name:  "empty",
							value: JSON_VALUE_ARRAY,
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
			tokens: parseTestCases[i].tokens,
		}
		parser.parse()

		if !equal(&parser.ast, &parseTestCases[i].ast) {
			t.Fatal("unexpected AST")
		}
	}
}
