#include <stdio.h>
#include <stdlib.h>

#include "util.h"

void log_and_exit(const char* message) {
    fprintf(stderr, "%s", message);
    exit(1);
}
