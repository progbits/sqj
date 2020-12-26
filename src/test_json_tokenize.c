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
        LEFT_CURLY_BRACKET,
        STRING,
        COLON,
        NUMBER,
        COMMA,
        STRING,
        COLON,
        NUMBER,
        COMMA,
        STRING,
        COLON,
        TRUE,
        COMMA,
        STRING,
        COLON,
        FALSE,
        COMMA,
        STRING,
        COLON,
        STRING,
        RIGHT_CURLY_BRACKET,
    };
    assert_test_case(test_0, expected_0, sizeof(expected_0) / sizeof(expected_0[0]));

    const char* test_1 = "{\"foo\": [{\"bar\": 1}, {\"baz\": 2}]}";
    const JSON_TOKEN expected_1[] = {
            LEFT_CURLY_BRACKET,
            STRING,
            COLON,
            LEFT_SQUARE_BRACKET,
            LEFT_CURLY_BRACKET,
            STRING,
            COLON,
            NUMBER,
            RIGHT_CURLY_BRACKET,
            COMMA,
            LEFT_CURLY_BRACKET,
            STRING,
            COLON,
            NUMBER,
            RIGHT_CURLY_BRACKET,
            RIGHT_SQUARE_BRACKET,
            RIGHT_CURLY_BRACKET,
    };
    assert_test_case(test_1, expected_1, sizeof(expected_1) / sizeof(expected_1[0]));
}