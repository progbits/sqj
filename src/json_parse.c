#include "json_parse.h"
#include "util.h"

void parse_object(JSONNode *root, Tokens *tokens, size_t offset);

Token *peek_token(Tokens *tokens, size_t offset) {
    if (offset >= tokens->size) {
        return NULL;
    }
    return tokens->data[offset];
}


// Parse an array.
void parse_array(JSONNode *node, Tokens *tokens, size_t offset) {
    node->value = JSON_VALUE_ARRAY;

    // Exit early for empty arrays.
    if (tokens->data[offset]->type == JSON_TOKEN_RIGHT_SQUARE_BRACKET) {
        return;
    }

    // Parse the values of the array.
    while (offset < tokens->size) {
        // Allocate a new value.
        JSONNode *value = calloc(1, sizeof(JSONNode));

        // Parse the value.
        switch (tokens->data[offset]->type) {
            case (JSON_TOKEN_FALSE): {
                value->value = JSON_VALUE_FALSE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_NULL): {
                value->value = JSON_VALUE_FALSE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_TRUE): {
                value->value = JSON_VALUE_TRUE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                value->value = JSON_VALUE_OBJECT;
                parse_object(value, tokens, offset + 1);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                value->value = JSON_VALUE_ARRAY;
                parse_array(value, tokens, offset + 1);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                value->value = JSON_VALUE_NUMBER;
                value->number_value = tokens->data[offset]->number;
                ++offset;
                break;
            }
            case (JSON_TOKEN_STRING): {
                value->value = JSON_VALUE_STRING;
                value->string_value = strdup(tokens->data[offset]->string);
                ++offset;
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        node->next = value;
        node = node->next;

        // Parse the next member.
        if (tokens->data[offset]->type == JSON_TOKEN_COMMA) {
            ++offset;
            continue;
        }

        // End of array.
        if (tokens->data[offset]->type == JSON_TOKEN_RIGHT_SQUARE_BRACKET) {
            ++offset;
            break;
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse an object.
void parse_object(JSONNode *node, Tokens *tokens, size_t offset) {
    node->value = JSON_VALUE_OBJECT;

    // Exit early for empty objects.
    if (tokens->data[offset]->type == JSON_TOKEN_RIGHT_CURLY_BRACKET) {
        return;
    }

    // Parse the members of the object.
    while (offset < tokens->size) {
        // Object has at least one member.
        if (tokens->data[offset]->type != JSON_TOKEN_STRING) {
            log_and_exit("expected a value of type string\n");
        }

        // Parse the member name.
        JSONNode *member = calloc(1, sizeof(JSONNode));
        member->name = strdup(tokens->data[offset]->string);
        ++offset;

        if (tokens->data[offset]->type != JSON_TOKEN_COLON) {
            log_and_exit("expected JSON_TOKEN_COLON\n");
        }
        ++offset;

        // Parse the member value.
        switch (tokens->data[offset]->type) {
            case (JSON_TOKEN_FALSE): {
                member->value = JSON_VALUE_FALSE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_NULL): {
                member->value = JSON_VALUE_FALSE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_TRUE): {
                member->value = JSON_VALUE_TRUE;
                ++offset;
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                member->value = JSON_VALUE_OBJECT;
                parse_object(member, tokens, offset + 1);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                member->value = JSON_VALUE_ARRAY;
                parse_array(member, tokens, offset + 1);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                member->value = JSON_VALUE_NUMBER;
                member->number_value = tokens->data[offset]->number;
                ++offset;
                break;
            }
            case (JSON_TOKEN_STRING): {
                member->value = JSON_VALUE_STRING;
                member->string_value = strdup(tokens->data[offset]->string);
                ++offset;
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        node->next = member;
        node = node->next;

        // Parse the next member.
        if (tokens->data[offset]->type == JSON_TOKEN_COMMA) {
            ++offset;
            continue;
        }

        // End of object.
        if (tokens->data[offset]->type == JSON_TOKEN_RIGHT_CURLY_BRACKET) {
            ++offset;
            break;
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse a JSON AST from a stream of tokens.
void parse(Tokens *tokens, JSONNode **ast) {
    // At the moment, we only consider our root node can be of type Object.
    *ast = calloc(1, sizeof(JSONNode));
    switch (tokens->data[0]->type) {
        case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
            parse_object(*ast, tokens, 1);
            break;
        }
        case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
            parse_array(*ast, tokens, 1);
            break;
        }
        default: {
            log_and_exit("expected first token to be { or [\n");
        }
    }
}
