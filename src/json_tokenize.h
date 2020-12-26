#ifndef JSON_TOKENIZE_H
#define JSON_TOKENIZE_H

#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

#include "vector.h"

typedef enum JSON_TOKEN {
    LEFT_SQUARE_BRACKET,
    LEFT_CURLY_BRACKET,
    RIGHT_SQUARE_BRACKET,
    RIGHT_CURLY_BRACKET,
    COLON,
    COMMA,
    WHITESPACE,
    FALSE,
    JSON_NULL,
    TRUE,
    OBJECT,
    ARRAY,
    NUMBER,
    STRING
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
