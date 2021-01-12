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
    // Count how many ' we need to escape.
    int n_esc = 0;
    for (int i = 0; i < strlen(str); i++) {
        if (str[i] == '\'') {
            ++n_esc;
        }
    }

    // Allocate space for leading and trailing ', plus any additional ' that we
    // need to escape.
    char* escaped = calloc(strlen(str) + n_esc + 3, sizeof(char));
    escaped[0] = '\'';
    for (int i = 0, j = 1; i < strlen(str); i++, j++) {
        if (str[i] == '\'') {
            escaped[j++] = '\'';
            escaped[j] = str[i];
            continue;
        }
        escaped[j] = str[i];
    }
    escaped[strlen(str) + 1] = '\'';

    return escaped;
}
