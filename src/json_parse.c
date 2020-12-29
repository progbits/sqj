#include <stddef.h>

#include "json_parse.h"
#include "util.h"

void parse_object(Token** tokens, JSONNode* root);

// Parse an array.
//
// When invoked, *tokens* should point to the first token after the opening
// square brace of the array. After the array has been parsed, *tokens* will
// be left pointing to the next token after the closing square brace of the
// array.
//
// - *tokens* - The stream of tokens from which to parse the array.
// - *node*   - An AST node of JSON_VALUE_ARRAY type. This nodes
//              *node->values* variable will hold the values of the array.
//
void parse_array(Token** tokens, JSONNode* node) {
    size_t n_values = 0;
    while (tokens) {
        node->values =
            realloc(node->values, (n_values + 1) * sizeof(struct JSONNode));

        JSONNode* value = &(node->values[n_values]);
        memset(value, 0, sizeof(struct JSONNode));

        switch ((*tokens)->type) {
            case (JSON_TOKEN_RIGHT_SQUARE_BRACKET): {
                ++(*tokens);
                return; // Empty array.
            }
            case (JSON_TOKEN_FALSE): {
                value->value = JSON_VALUE_FALSE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_NULL): {
                value->value = JSON_VALUE_FALSE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_TRUE): {
                value->value = JSON_VALUE_TRUE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                value->value = JSON_VALUE_OBJECT;
                ++(*tokens);
                parse_object(tokens, value);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                value->value = JSON_VALUE_ARRAY;
                ++(*tokens);
                parse_array(tokens, value);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                value->value = JSON_VALUE_NUMBER;
                value->number_value = (*tokens)->number;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_STRING): {
                value->value = JSON_VALUE_STRING;
                value->string_value = strdup((*tokens)->string);
                ++(*tokens);
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        if ((*tokens)->type == JSON_TOKEN_COMMA) {
            ++n_values;
            ++(*tokens);
            continue; // Parse the next array value.
        } else if ((*tokens)->type == JSON_TOKEN_RIGHT_SQUARE_BRACKET) {
            ++(*tokens);
            break; // End of array.
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse an object from a stream of tokens.
//
// When invoked, *tokens* should point to the first token after the opening
// curly brace of the object.. After the object has been parsed, *tokens* will
// be left pointing to the next token after the closing curly brace of the
// object.
//
// - *tokens* - The stream of tokens from which to parse the object.
// - *node*   - An AST node of JSON_VALUE_OBJECT type. This nodes
//              *node->members* variable will hold the members of the object.
//
void parse_object(Token** tokens, JSONNode* node) {
    size_t n_members = 0;
    while (*tokens) {
        node->members =
            realloc(node->members, (n_members + 1) * sizeof(struct JSONNode));

        JSONNode* member = &(node->members[n_members]);
        memset(member, 0, sizeof(struct JSONNode));

        if ((*tokens)->type != JSON_TOKEN_STRING) {
            log_and_exit("expected a value of type string\n");
        }
        member->name = strdup((*tokens)->string);

        if ((++(*tokens))->type != JSON_TOKEN_COLON) {
            log_and_exit("expected JSON_TOKEN_COLON\n");
        }

        // Parse the member value.
        switch ((++(*tokens))->type) {
            case (JSON_TOKEN_RIGHT_CURLY_BRACKET): {
                ++(*tokens);
                return; // Empty object.
            }
            case (JSON_TOKEN_FALSE): {
                member->value = JSON_VALUE_FALSE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_NULL): {
                member->value = JSON_VALUE_FALSE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_TRUE): {
                member->value = JSON_VALUE_TRUE;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_LEFT_CURLY_BRACKET): {
                member->value = JSON_VALUE_OBJECT;
                ++(*tokens);
                parse_object(tokens, member);
                break;
            }
            case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
                member->value = JSON_VALUE_ARRAY;
                parse_array(tokens, member);
                break;
            }
            case (JSON_TOKEN_NUMBER): {
                member->value = JSON_VALUE_NUMBER;
                member->number_value = (*tokens)->number;
                ++(*tokens);
                break;
            }
            case (JSON_TOKEN_STRING): {
                member->value = JSON_VALUE_STRING;
                member->string_value = strdup((*tokens)->string);
                ++(*tokens);
                break;
            }
            default: {
                log_and_exit("unexpected token\n");
            }
        }

        if ((*tokens)->type == JSON_TOKEN_COMMA) {
            ++n_members;
            ++(*tokens);
            continue; // Parse next member.
        } else if ((*tokens)->type == JSON_TOKEN_RIGHT_CURLY_BRACKET) {
            ++(*tokens);
            break; // End of object.
        }
        log_and_exit("unexpected token\n");
    }
}

// Parse a JSON AST from a stream of tokens.
//
// *tokens* - Stream of tokens to be parsed.
// *ast* - Assigned to the root of the AST.
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
            (*ast)->value = JSON_VALUE_OBJECT;
            ++tokens;
            parse_object(&tokens, *ast);
            break;
        }
        case (JSON_TOKEN_LEFT_SQUARE_BRACKET): {
            (*ast)->value = JSON_VALUE_ARRAY;
            ++tokens;
            parse_array(&tokens, *ast);
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
