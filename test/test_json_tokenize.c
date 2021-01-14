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

START_TEST(test_json_tokenize_2) {
    const char* raw_json =
        "{\"foo\": \"\\\"\\\"\", \"bar\": \"\\\"hello, world\\\"\"}";

    Token* tokens = NULL;
    size_t n_tokens = 0;
    tokenize(raw_json, &tokens, &n_tokens);

    const JSON_TOKEN expected_tokens[] = {
        JSON_TOKEN_LEFT_CURLY_BRACKET,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_STRING,
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

    // Clean up.
    free(tokens);
}
END_TEST

START_TEST(test_json_tokenize_3) {
    const char* raw_json = "[{\"empty\": []}]";

    Token* tokens = NULL;
    size_t n_tokens = 0;
    tokenize(raw_json, &tokens, &n_tokens);

    const JSON_TOKEN expected_tokens[] = {
        JSON_TOKEN_LEFT_SQUARE_BRACKET,
        JSON_TOKEN_LEFT_CURLY_BRACKET,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_LEFT_SQUARE_BRACKET,
        JSON_TOKEN_RIGHT_SQUARE_BRACKET,
        JSON_TOKEN_RIGHT_CURLY_BRACKET,
        JSON_TOKEN_RIGHT_SQUARE_BRACKET,
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

START_TEST(test_json_tokenize_4) {
    const char* raw_json =
        "{"
        "   \"α\": 0.0072973525693,"
        "   \"γ\": 0.5772156649015328606065120900824024310421,"
        "   \"δ\": 4.669201609102990671853203820466,"
        "   \"ϵ\": 8.854187812813e12,"
        "   \"ζ\": 1.202056903159594285399738161511449990764986292,"
        "   \"θ\": 90,"
        "   \"μ\": 1.2566370614E-6,"
        "   \"ψ\": 3.359885666243177553172011302918927179688905133732"
        "}";

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
        JSON_TOKEN_NUMBER,
        JSON_TOKEN_COMMA,
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
        JSON_TOKEN_NUMBER,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_NUMBER,
        JSON_TOKEN_COMMA,
        JSON_TOKEN_STRING,
        JSON_TOKEN_COLON,
        JSON_TOKEN_NUMBER,
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
    tcase_add_test(tc_core, test_json_tokenize_2);
    tcase_add_test(tc_core, test_json_tokenize_3);
    tcase_add_test(tc_core, test_json_tokenize_4);

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
