#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "json_tokenize.h"

void print_json_token(Token* token) {
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

int is_json_whitespace(const char* value) {
    switch (*value) {
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
int starts_with(const char* prefix, const char* string) {
    if (strlen(prefix) > strlen(string)) {
        return 0;
    }

    return memcmp(prefix, string, strlen(prefix)) == 0;
}

void tokenize(const char* input, Token** tokens, size_t* n_tokens) {
    size_t allocated = 0;
    while (*input != '\0') {
        // Allocate some space for the next token.
        *tokens = realloc(*tokens, (allocated + 1) * sizeof(Token));
        Token* token = &((*tokens)[allocated++]);

        // Skip whitespace.
        while (is_json_whitespace(input)) {
            if (*(++input) == '\0') {
                return;
            }
        }

        // Handle structural characters.
        switch (*input) {
            case '[': {
                token->type = JSON_TOKEN_LEFT_SQUARE_BRACKET;
                ++input;
                continue;
            }
            case '{': {
                token->type = JSON_TOKEN_LEFT_CURLY_BRACKET;
                ++input;
                continue;
            }
            case ']': {
                token->type = JSON_TOKEN_RIGHT_SQUARE_BRACKET;
                ++input;
                continue;
            }
            case '}': {
                token->type = JSON_TOKEN_RIGHT_CURLY_BRACKET;
                ++input;
                continue;
            }
            case ':': {
                token->type = JSON_TOKEN_COLON;
                ++input;
                continue;
            }
            case ',': {
                token->type = JSON_TOKEN_COMMA;
                ++input;
                continue;
            }
        }

        // Handle boolean literals.
        if (starts_with("true", input)) {
            token->type = JSON_TOKEN_TRUE;
            input += strlen("true");
            continue;
        }
        if (starts_with("false", input)) {
            token->type = JSON_TOKEN_FALSE;
            input += strlen("false");
            continue;
        }

        // Handle null literal.
        if (starts_with("null", input)) {
            token->type = JSON_TOKEN_NULL;
            input += strlen("null");
            continue;
        }

        // Consume numeric literals.
        if ((*input >= '0' && *input <= '9') || *input == '-') {
            char* end;
            double value = strtod(input, &end);
            token->type = JSON_TOKEN_NUMBER;
            token->number = value;
            input = end;
            continue;
        }

        // Must be consuming a string.
        if (*input == '\"') {
            const char* start = ++input;
            while (input && *input != '\"') {
                ++input;
            }
            token->type = JSON_TOKEN_STRING;
            token->string = strndup(start, input - start);
            ++input;
            continue;
        }

        // This should never happen.
        fprintf(stderr, "tokenization failed!\n");
        exit(1);
    }
    *n_tokens = allocated;
}
