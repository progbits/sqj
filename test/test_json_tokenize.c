#include <check.h>
#include <stddef.h>
#include <stdlib.h>

#include "../src/json_tokenize.h"

START_TEST(test_json_tokenize_0) {
    // Arrange.
    const char* raw_json = "{\"foo\": 1, \"bar\": 3.14, \"baz\": true, "
                           "\"foobar\": false, \"foobaz\": \"hello, world\"}";
    Token* tokens = NULL;
    size_t n_tokens = 0;
    tokenize(raw_json, &tokens, &n_tokens);

    const JSON_TOKEN expected_tokens[] = {
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

    // Assert.
    ck_assert_int_eq(n_tokens,
                     sizeof(expected_tokens) / sizeof(expected_tokens[0]));
    for (size_t i = 0; i < n_tokens; i++) {
        ck_assert_int_eq(tokens[i].type, expected_tokens[i]);
    }
    free(tokens);
}
END_TEST

START_TEST(test_json_tokenize_1) {
    const char* raw_json = "{\"foo\": [{\"bar\": 1}, {\"baz\": 2}]}";

    Token* tokens = NULL;
    size_t n_tokens = 0;
    tokenize(raw_json, &tokens, &n_tokens);

    const JSON_TOKEN expected_tokens[] = {
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

    // Assert.
    ck_assert_int_eq(n_tokens,
                     sizeof(expected_tokens) / sizeof(expected_tokens[0]));
    for (size_t i = 0; i < n_tokens; i++) {
        ck_assert_int_eq(tokens[i].type, expected_tokens[i]);
    }

    // Clean up.
    free(tokens);
}
END_TEST

Suite* tokenizer_suite(void) {
    Suite* s = suite_create("json_tokenize");

    TCase* tc_core = tcase_create("Core");
    tcase_add_test(tc_core, test_json_tokenize_0);
    tcase_add_test(tc_core, test_json_tokenize_1);

    suite_add_tcase(s, tc_core);
    return s;
}

int main(void) {
    int number_failed;
    Suite* s;
    SRunner* sr;

    s = tokenizer_suite();
    sr = srunner_create(s);

    srunner_run_all(sr, CK_VERBOSE);
    number_failed = srunner_ntests_failed(sr);
    srunner_free(sr);
    return (number_failed == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}
