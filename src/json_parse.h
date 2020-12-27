#ifndef JSON_PARSE_H
#define JSON_PARSE_H

#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

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

    // The next member of an object or the next value of an array.
    struct JSONNode* next;

    // Object members.
    struct JSONNode* members;

    // Array items.
    struct JSONNode* values;

    // Value for tokens of type NUMBER.
    double number_value;

    // Value for tokens of type STRING.
    char* string_value;
} JSONNode;

// Parse a collection of JSON tokens into an AST.
void parse(Tokens* tokens, JSONNode** ast);

#endif // JSON_PARSE_H
