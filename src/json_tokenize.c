#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

#include "json_tokenize.h"

void init_tokens(Tokens *tokens, size_t size) {
    if (size == 0) {
        size = 1;
    }
    tokens->data = malloc(size * sizeof(Token *));
    tokens->size = 0;
    tokens->allocated_size = size;
}

void push_token(Tokens *tokens, Token *token) {
    if (tokens->size == tokens->allocated_size) {
        tokens->allocated_size *= 2;
        tokens->data = realloc(tokens->data, tokens->allocated_size * sizeof(Token *));
    }
    tokens->data[tokens->size++] = token;
}

void print_json_token(Token *token) {
    switch (token->type) {
        case JSON_TOKEN_LEFT_SQUARE_BRACKET: {
            printf("LEFT_SQUARE_BRACKET\n");
            break;
        }
        case JSON_TOKEN_LEFT_CURLY_BRACKET: {
            printf("LEFT_CURLY_BRACKET\n");
            break;
        }
        case JSON_TOKEN_RIGHT_SQUARE_BRACKET: {
            printf("RIGHT_SQUARE_BRACKET\n");
            break;
        }
        case JSON_TOKEN_RIGHT_CURLY_BRACKET: {
            printf("RIGHT_CURLY_BRACKET\n");
            break;
        }
        case JSON_TOKEN_COLON: {
            printf("COLON\n");
            break;
        }
        case JSON_TOKEN_COMMA: {
            printf("COMMA\n");
            break;
        }
        case JSON_TOKEN_WHITESPACE: {
            printf("WHITESPACE\n");
            break;
        }
        case JSON_TOKEN_FALSE: {
            printf("FALSE\n");
            break;
        }
        case JSON_TOKEN_NULL: {
            printf("JSON_NULL\n");
            break;
        }
        case JSON_TOKEN_TRUE: {
            printf("TRUE\n");
            break;
        }
        case JSON_TOKEN_OBJECT: {
            printf("OBJECT\n");
            break;
        }
        case JSON_TOKEN_ARRAY: {
            printf("ARRAY\n");
            break;
        }
        case JSON_TOKEN_NUMBER: {
            printf("NUMBER: %f\n", token->number);
            break;
        }
        case JSON_TOKEN_STRING: {
            printf("STRING: %s\n", token->string);
            break;
        }
        default: {
            assert("unknown token");
        }
    }
}

int is_json_whitespace(char value) {
    switch (value) {
        case '\x20':
        case '\x09':
        case '\x0A':
        case '\x0D':
            return 1;
        default:
            return 0;
    }
}

// Check if a string starts with a prefix.
int starts_with(const char *prefix, const char *string) {
    if (strlen(prefix) > strlen(string)) {
        return 0;
    }

    return memcmp(prefix, string, strlen(prefix)) == 0;
}

// Tokenize a JSON input.
//
// Tokenize the RCF7159 JSON grammar:
//      JSON-text = ws value ws.
//
// Returns the number of data.
void tokenize(const char *input, size_t size, Tokens *tokens) {
    init_tokens(tokens, 1);
    size_t index = 0;
    while (index < size) {
        // Allocate a new Token.
        Token *token = calloc(1, sizeof(Token));

        // Skip whitespace.
        while (is_json_whitespace(input[index])) {
            ++index;
        }

        // Handle structural characters.
        switch (input[index]) {
            case '[': {
                token->type = JSON_TOKEN_LEFT_SQUARE_BRACKET;
                push_token(tokens, token);
                ++index;
                continue;
            }
            case '{': {
                token->type = JSON_TOKEN_LEFT_CURLY_BRACKET;
                push_token(tokens, token);
                ++index;
                continue;
            }
            case ']': {
                token->type = JSON_TOKEN_RIGHT_SQUARE_BRACKET;
                push_token(tokens, token);
                ++index;
                continue;
            }
            case '}': {
                token->type = JSON_TOKEN_RIGHT_CURLY_BRACKET;
                push_token(tokens, token);
                ++index;
                continue;
            }
            case ':': {
                token->type = JSON_TOKEN_COLON;
                push_token(tokens, token);
                ++index;
                continue;
            }
            case ',': {
                token->type = JSON_TOKEN_COMMA;
                push_token(tokens, token);
                ++index;
                continue;
            }
        }

        // Handle boolean literals.
        if (starts_with("true", input + index)) {
            token->type = JSON_TOKEN_TRUE;
            push_token(tokens, token);
            index += strlen("true");
            continue;
        }
        if (starts_with("false", input + index)) {
            token->type = JSON_TOKEN_FALSE;
            push_token(tokens, token);
            index += strlen("false");
            continue;
        }

        // Handle null literal.
        if (starts_with("null", input + index)) {
            token->type = JSON_TOKEN_NULL;
            push_token(tokens, token);
            index += strlen("null");
            continue;
        }

        // Consume numeric literals.
        if ((input[index] >= '0' && input[index] <= '9') || input[index] == '-') {
            char *end;
            double value = strtod(input + index, &end);
            token->type = JSON_TOKEN_NUMBER;
            token->number = value;
            push_token(tokens, token);
            index += end - (input + index);
            continue;
        }

        // Must be consuming a string.
        if (input[index] == '\"') {
            size_t consume = ++index;
            while (consume < size && input[consume] != '\"') {
                ++consume;
            }
            token->type = JSON_TOKEN_STRING;
            token->string = calloc(((consume - index) + 1), sizeof(char));
            memcpy(token->string, input + index, consume - index);
            push_token(tokens, token);
            index = ++consume;
            continue;
        }

        // This should never happen.
        fprintf(stderr, "tokenization failed!\n");
        exit(1);
    }

    for (int i = 0; i < tokens->size; i++) {
        print_json_token(tokens->data[i]);
    }
}
