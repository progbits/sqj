#ifndef JSON_PARSE_H
#define JSON_PARSE_H

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "json_tokenize.h"

// Allowed JSON values.
//
// RFC 7159 - Section 2.
typedef enum JSON_VALUE {
    JSON_VALUE_OBJECT,
    JSON_VALUE_ARRAY,
    JSON_VALUE_NUMBER,
    JSON_VALUE_STRING,
    JSON_VALUE_NULL,
    JSON_VALUE_TRUE,
    JSON_VALUE_FALSE
} JSON_VALUE;

// JSONNode represents a node the tree representation of a JSON string.
typedef struct JSONNode {
    // The value type of this node.
    JSON_VALUE value;

    // Name of object member.
    char* name;

    // Object members.
    struct JSONNode* members;

    // Number of object members.
    size_t n_members;

    // Array items.
    struct JSONNode* values;

    // Number of array values.
    size_t n_values;

    // Value for tokens of type NUMBER.
    double number_value;

    // Value for tokens of type STRING.
    char* string_value;
} JSONNode;

// Parse a collection of JSON tokens into an AST.
void parse(Token* tokens, JSONNode** ast);

// Clean up memory allocated for an AST.
void delete_ast(JSONNode* ast);

// Pretty print a JSON AST to a file stream.
//
// If compact != 0 then all formatting related whitespace will be omitted from
// the output.
void pretty_print(JSONNode* ast, FILE* stream, int compact);

#endif // JSON_PARSE_H
