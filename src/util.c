#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "util.h"

void log_and_exit(const char* format, ...) {
    va_list arglist;
    va_start(arglist, format);
    vfprintf(stderr, format, arglist);
    va_end(arglist);
    exit(EXIT_FAILURE);
}

char* escape_string(const char* str) {
    if (strlen(str) == 0) {
        return NULL;
    }

    // Since we are enclosing our string in single quotation marks ('), we need
    // to escape any exiting ' characters. n_esc starts at 2 to account for the
    // leading and trailing single quotation marks.
    int n_esc = 2;
    for (int i = 0; i < strlen(str); i++) {
        if (str[i] == '\'') {
            ++n_esc;
        }
    }

    // Allocate space for the escaped string.
    char* escaped = calloc(strlen(str) + n_esc + 1, sizeof(char));

    // Build the escaped string.
    int j = 0;
    escaped[j++] = '\'';
    for (int i = 0; i < strlen(str); i++, j++) {
        if (str[i] == '\'') {
            escaped[j++] = '\'';
            escaped[j] = str[i];
            continue;
        }
        escaped[j] = str[i];
    }
    escaped[j] = '\'';

    return escaped;
}

char* unescape_string(const char* str) {
    if (strlen(str) == 0) {
        return NULL;
    }

    // Count how many ' characters we escaped.
    int n_esc = 0;
    for (int i = 0; i < strlen(str); i++) {
        if (str[i] == '\'') {
            ++n_esc;
            ++i;
        }
    }

    // Allocate space into which to retrieve the original string.
    char* original = calloc((strlen(str) - n_esc) + 1, sizeof(char));

    // Unescape string.
    for (int i = 0, j = 1; j < strlen(str) - 1; i++, j++) {
        if (str[j] == '\'') {
            original[i] = str[j++];
            continue;
        }
        original[i] = str[j];
    }
    return original;
}
