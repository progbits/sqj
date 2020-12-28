#include <check.h>
#include <stdlib.h>

#include "../src/json_parse.h"

START_TEST(test_json_parse_simple_string) {
    // Arrange.
    JSONNode* ast = NULL;

    const Token tokens[] = {
        {.type = JSON_TOKEN_STRING, .string = "hello, world"}};

    // Act.
    parse((Token*)tokens, &ast);

    // Assert.
    ck_assert_int_eq(ast->value, JSON_VALUE_STRING);
    ck_assert_str_eq(ast->string_value, "hello, world");

    // Clean up.
    free(ast);
}
END_TEST

START_TEST(test_json_parse_simple_number) {
    // Arrange.
    JSONNode* ast = NULL;

    const Token tokens[] = {{.type = JSON_TOKEN_NUMBER, .number = 3.14}};

    // Act.
    parse((Token*)tokens, &ast);

    // Assert.
    ck_assert_int_eq(ast->value, JSON_VALUE_NUMBER);
    ck_assert_double_eq(ast->number_value, 3.14);

    // Clean up.
    free(ast);
}
END_TEST

START_TEST(test_json_parse_simple_null) {
    // Arrange.
    JSONNode* ast = NULL;

    const Token tokens[] = {{.type = JSON_TOKEN_NULL}};

    // Act.
    parse((Token*)tokens, &ast);

    // Assert.
    ck_assert_int_eq(ast->value, JSON_VALUE_NULL);

    // Clean up.
    free(ast);
}
END_TEST

START_TEST(test_json_parse_simple_true) {
    // Arrange.
    JSONNode* ast = NULL;

    const Token tokens[] = {{.type = JSON_TOKEN_TRUE}};

    // Act.
    parse((Token*)tokens, &ast);

    // Assert.
    ck_assert_int_eq(ast->value, JSON_VALUE_TRUE);

    // Clean up.
    free(ast);
}
END_TEST

START_TEST(test_json_parse_simple_false) {
    // Arrange.
    JSONNode* ast = NULL;

    const Token tokens[] = {{.type = JSON_TOKEN_FALSE}};

    // Act.
    parse((Token*)tokens, &ast);

    // Assert.
    ck_assert_int_eq(ast->value, JSON_VALUE_FALSE);

    // Clean up.
    free(ast);
}
END_TEST

Suite* parser_suite(void) {
    Suite* s = suite_create("json_parse");

    TCase* tc_core = tcase_create("Core");
    tcase_add_test(tc_core, test_json_parse_simple_string);
    tcase_add_test(tc_core, test_json_parse_simple_number);
    tcase_add_test(tc_core, test_json_parse_simple_null);
    tcase_add_test(tc_core, test_json_parse_simple_true);
    tcase_add_test(tc_core, test_json_parse_simple_false);

    suite_add_tcase(s, tc_core);
    return s;
}

int main(void) {
    int number_failed;
    Suite* s;
    SRunner* sr;

    s = parser_suite();
    sr = srunner_create(s);

    srunner_run_all(sr, CK_VERBOSE);
    number_failed = srunner_ntests_failed(sr);
    srunner_free(sr);
    return (number_failed == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}
