#ifndef JSON_TOKENIZE_H
#define JSON_TOKENIZE_H

#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

#include "vector.h"

typedef enum JSON_TOKEN {
    JSON_TOKEN_LEFT_SQUARE_BRACKET,
    JSON_TOKEN_LEFT_CURLY_BRACKET,
    JSON_TOKEN_RIGHT_SQUARE_BRACKET,
    JSON_TOKEN_RIGHT_CURLY_BRACKET,
    JSON_TOKEN_COLON,
    JSON_TOKEN_COMMA,
    JSON_TOKEN_WHITESPACE,
    JSON_TOKEN_FALSE,
    JSON_TOKEN_NULL,
    JSON_TOKEN_TRUE,
    JSON_TOKEN_OBJECT,
    JSON_TOKEN_ARRAY,
    JSON_TOKEN_NUMBER,
    JSON_TOKEN_STRING
} JSON_TOKEN;

typedef struct JSONToken {
    enum JSON_TOKEN type;
    void* value;
} JSONToken;

// Tokenize a JSON input.
//
// Tokenize the RCF7159 JSON grammar.
//
// JSON-test = ws value ws.
void tokenize(const char* input, size_t size, Vec** tokens);

#endif // JSON_TOKENIZE_H
