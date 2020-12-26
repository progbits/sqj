#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

#include "json_tokenize.h"
#include "vector.h"

void print_json_token(JSON_TOKEN token) {
    switch (token) {
        case LEFT_SQUARE_BRACKET: {
            printf("LEFT_SQUARE_BRACKET\n");
            break;
        }
        case LEFT_CURLY_BRACKET: {
            printf("LEFT_CURLY_BRACKET\n");
            break;
        }
        case RIGHT_SQUARE_BRACKET: {
            printf("RIGHT_SQUARE_BRACKET\n");
            break;
        }
        case RIGHT_CURLY_BRACKET: {
            printf("RIGHT_CURLY_BRACKET\n");
            break;
        }
        case COLON: {
            printf("COLON\n");
            break;
        }
        case COMMA: {
            printf("COMMA\n");
            break;
        }
        case WHITESPACE: {
            printf("WHITESPACE\n");
            break;
        }
        case FALSE: {
            printf("FALSE\n");
            break;
        }
        case JSON_NULL: {
            printf("JSON_NULL\n");
            break;
        }
        case TRUE: {
            printf("TRUE\n");
            break;
        }
        case OBJECT: {
            printf("OBJECT\n");
            break;
        }
        case ARRAY: {
            printf("ARRAY\n");
            break;
        }
        case NUMBER: {
            printf("NUMBER\n");
            break;
        }
        case STRING: {
            printf("STRING\n");
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
// Tokenize the RCF7159 JSON grammar.
//
// JSON-text = ws value ws.
void tokenize(const char* input, size_t size, Vec** tokens) {
    *tokens = malloc(sizeof (Vec));
    new_vec(*tokens, 0);
    size_t index = 0;
    while (index < size) {
        // Skip whitespace.
        while (is_json_whitespace(input[index])) {
            ++index;
        }

        // Handle structural characters.
        switch (input[index]) {
            case '[': {
                vec_push_back(*tokens, LEFT_SQUARE_BRACKET);
                ++index;
                continue;
            }
            case '{': {
                vec_push_back(*tokens, LEFT_CURLY_BRACKET);
                ++index;
                continue;
            }
            case ']': {
                vec_push_back(*tokens, RIGHT_SQUARE_BRACKET);
                ++index;
                continue;
            }
            case '}': {
                vec_push_back(*tokens, RIGHT_CURLY_BRACKET);
                ++index;
                continue;
            }
            case ':': {
                vec_push_back(*tokens, COLON);
                ++index;
                continue;
            }
            case ',': {
                vec_push_back(*tokens, COMMA);
                ++index;
                continue;
            }
        }

        // Handle boolean literals.
        if (starts_with("true", input + index)) {
            vec_push_back(*tokens, TRUE);
            input += strlen("true");
            continue;
        }
        if (starts_with("false", input + index)) {
            vec_push_back(*tokens, FALSE);
            input += strlen("false");
            continue;
        }

        // Handle null literal.
        if (starts_with("null", input + index)) {
            vec_push_back(*tokens, JSON_NULL);
            input += strlen("null");
            continue;
        }

        // Consume numeric literals.
        if ((input[index] >= '0' && input[index] <= '9') || input[index] == '-') {
            char* end;
            double value = strtod(input + index, &end);
            printf("Parsed numeric literal %f\n", value);
            vec_push_back(*tokens, NUMBER);
            index += end - (input + index);
            continue;
        }

        // Must be consuming a string.
        if (input[index] == '\"') {
            size_t consume = ++index;
            while (consume < size && input[consume] != '\"') {
                ++consume;
            }
            vec_push_back(*tokens, STRING);
            index = ++consume;
            continue;
        }

        // This should never happen.
        assert("Tokenization failed!");
        break;
    }

    for (int i = 0; i < (*tokens)->occupied; ++i) {
        print_json_token((*tokens)->data[i]);
    }
}
