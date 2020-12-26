#include <stdlib.h>

#include "vector.h"

void new_vec(Vec* vec, size_t size) {
    if (size == 0) {
        size = 1;
    }
    vec->data = malloc(size * sizeof(int));
    vec->occupied = 0;
    vec->actual = size;
}

void vec_push_back(Vec* vec, int value) {
    if (vec->occupied == vec->actual) {
        vec->actual *= 2;
        vec->data = realloc(vec->data, vec->actual * sizeof(int));
    }
    vec->data[vec->occupied++] = value;
}

void vec_destroy(Vec* vec) {
    free(vec->data);
    vec->data = NULL;
    vec->actual = 0;
    vec->occupied = 0;
}

