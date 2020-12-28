#include <stddef.h>

#include "json_parse.h"
#include "util.h"

void parse_object(JSONNode* root, Token* tokens);

// Parse an array.
//
// The parameter *node* points to a pre-allocated node for the array.
void parse_array(JSONNode* node, Token* tokens) {
    node->value = JSON_VALUE_ARRAY;
    while (tokens) {
        JSONNode* value = calloc(1, sizeof(JSONNode));
        switch (tokens->type) {
            case (JSON_TOKEN_RIGHT_SQUARE_BRACKET): {
                return; // Empty array.
            }
            case (JSON_TOKEN_FALSE): {
                value->value = JSON_VALUE_FALSE;
                break;
            }
            case (JSON_TOKEN_NULL): {
                value->value = JSON_VALUE_FALSE;
                break;
            }
            case (JSON_TOKEN_TRUE): {
                value->value = JSON_VALUE_TRUE;
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                value->value = JSON_VALUE_OBJECT;
                parse_object(value, tokens);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                value->value = JSON_VALUE_ARRAY;
                parse_array(value, tokens);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                value->value = JSON_VALUE_NUMBER;
                value->number_value = tokens->number;
                break;
            }
            case (JSON_TOKEN_STRING): {
                value->value = JSON_VALUE_STRING;
                value->string_value = strdup(tokens->string);
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        node->next = value;
        node = node->next;

        if (++tokens == NULL) {
            log_and_exit("unexpected end of token stream\n");
        }

        if (tokens->type == JSON_TOKEN_COMMA) {
            ++tokens;
            continue; // Parse the next array value.
        } else if (tokens->type == JSON_TOKEN_RIGHT_SQUARE_BRACKET) {
            ++tokens;
            break; // End of array.
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse an object.
void parse_object(JSONNode* node, Token* tokens) {
    node->value = JSON_VALUE_OBJECT;
    while (tokens) {
        if (tokens->type != JSON_TOKEN_STRING) {
            log_and_exit("expected a value of type string\n");
        }

        // Parse the member name.
        JSONNode* member = calloc(1, sizeof(JSONNode));
        member->name = strdup(tokens->string);

        if ((++tokens)->type != JSON_TOKEN_COLON) {
            log_and_exit("expected JSON_TOKEN_COLON\n");
        }

        // Parse the member value.
        switch ((++tokens)->type) {
            case (JSON_TOKEN_RIGHT_CURLY_BRACKET): {
                return; // Empty object.
            }
            case (JSON_TOKEN_FALSE): {
                member->value = JSON_VALUE_FALSE;
                break;
            }
            case (JSON_TOKEN_NULL): {
                member->value = JSON_VALUE_FALSE;
                break;
            }
            case (JSON_TOKEN_TRUE): {
                member->value = JSON_VALUE_TRUE;
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                member->value = JSON_VALUE_OBJECT;
                parse_object(member, tokens);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                member->value = JSON_VALUE_ARRAY;
                parse_array(member, tokens);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                member->value = JSON_VALUE_NUMBER;
                member->number_value = tokens->number;
                break;
            }
            case (JSON_TOKEN_STRING): {
                member->value = JSON_VALUE_STRING;
                member->string_value = strdup(tokens->string);
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        node->next = member;
        node = node->next;

        // Parse the next member.
        if (++tokens == NULL) {
            log_and_exit("unexpected end of token stream\n");
        }

        if (tokens->type == JSON_TOKEN_COMMA) {
            continue; // Parse next member.
        } else if (tokens->type == JSON_TOKEN_RIGHT_CURLY_BRACKET) {
            break; // End of object.
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse a JSON AST from a stream of tokens.
void parse(Token* tokens, JSONNode** ast) {
    // At the moment, we only consider our root node can be of type Object.
    *ast = calloc(1, sizeof(JSONNode));
    switch (tokens->type) {
        case (JSON_TOKEN_FALSE): {
            (*ast)->value = JSON_VALUE_FALSE;
            return;
        }
        case (JSON_TOKEN_NULL): {
            (*ast)->value = JSON_VALUE_NULL;
            return;
        }
        case (JSON_TOKEN_TRUE): {
            (*ast)->value = JSON_VALUE_TRUE;
            return;
        }
        case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
            parse_object(*ast, ++tokens);
            break;
        }
        case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
            parse_array(*ast, ++tokens);
            break;
        }
        case (JSON_TOKEN_NUMBER): {
            (*ast)->value = JSON_VALUE_NUMBER;
            (*ast)->number_value = tokens->number;
            return;
        }
        case (JSON_TOKEN_STRING): {
            (*ast)->value = JSON_VALUE_STRING;
            (*ast)->string_value = strdup(tokens->string);
            return;
        }
        default: {
            log_and_exit("unexpected token");
        }
    }
}
