#ifndef VECTOR_H
#define VECTOR_H

//
// Bare-bones dynamic array
//

struct Vec {
    int* data;
    size_t occupied;
    size_t actual;
};
typedef struct Vec Vec;

void new_vec(Vec* vec, size_t size);

void vec_push_back(Vec* vec, int value);

void vec_destroy(Vec* vec);

#endif // VECTOR_H