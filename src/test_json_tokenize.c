#include "json_tokenize.h"

void assert_test_case(const char* json, const JSON_TOKEN* expected, size_t expected_size) {
    // Arrange.
    Vec* tokens = malloc(sizeof(Vec));

    // Act.
    tokenize(json, strlen(json), &tokens);

    // Assert.
    assert(tokens->occupied == expected_size);
    for (size_t i = 0; i < expected_size; i++) {
        assert(tokens->data[i] == expected[i]);
    }

    // Clean up.
    vec_destroy(tokens);
    free(tokens);
}

int main() {
    const char* test_0 = "{\"foo\": 1, \"bar\": 3.14, \"baz\": true, \"foobar\": false, \"foobaz\": \"hello, world\"}";
    const JSON_TOKEN expected_0[] = {
        JSON_TOKEN_LEFT_CURLY_BRACKET,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_NUMBER,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_NUMBER,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_TRUE,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_FALSE,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_STRING,
        JSON_TOKEN_RIGHT_CURLY_BRACKET,
    };
    assert_test_case(test_0, expected_0, sizeof(expected_0) / sizeof(expected_0[0]));

    const char* test_1 = "{\"foo\": [{\"bar\": 1}, {\"baz\": 2}]}";
    const JSON_TOKEN expected_1[] = {
            JSON_TOKEN_LEFT_CURLY_BRACKET,
            JSON_TOKEN_STRING,
            JSON_TOKEN_COLON,
            JSON_TOKEN_LEFT_SQUARE_BRACKET,
            JSON_TOKEN_LEFT_CURLY_BRACKET,
            JSON_TOKEN_STRING,
            JSON_TOKEN_COLON,
            JSON_TOKEN_NUMBER,
            JSON_TOKEN_RIGHT_CURLY_BRACKET,
            JSON_TOKEN_COMMA,
            JSON_TOKEN_LEFT_CURLY_BRACKET,
            JSON_TOKEN_STRING,
            JSON_TOKEN_COLON,
            JSON_TOKEN_NUMBER,
            JSON_TOKEN_RIGHT_CURLY_BRACKET,
            JSON_TOKEN_RIGHT_SQUARE_BRACKET,
            JSON_TOKEN_RIGHT_CURLY_BRACKET,
    };
    assert_test_case(test_1, expected_1, sizeof(expected_1) / sizeof(expected_1[0]));
}