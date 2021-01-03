#include <stdio.h>
#include <stdlib.h>

#include "util.h"

void log_and_exit(const char* format, ...) {
    va_list arglist;
    va_start(arglist, format);
    vfprintf(stderr, format, arglist);
    va_end(arglist);
    exit(EXIT_FAILURE);
}
