#include <malloc.h>
#include <stdio.h>
#include <stdarg.h>

#include "../include/permutation.h"


// The initial permutation - IP.
// 
// As of now, it is hard-coded.
static uint8_t IP[] = 
{
    58, 50, 42, 34, 26, 18, 10, 2,
    60, 52, 44, 36, 28, 20, 12, 4,
    62, 54, 46, 38, 30, 22, 14, 6,
    64, 56, 48, 40, 32, 24, 16, 8,
    57, 49, 41, 33, 25, 17, 9, 1,
    59, 51, 43, 35, 27, 19, 11, 3,
    61, 53, 45, 37, 29, 21, 13, 5,
    63, 55, 47, 39, 31, 23, 15, 7
};


Perm *get_default_initial_permutation() {
    Perm *p = malloc(sizeof(Perm));
    p->size = 64;

    p->data = malloc(sizeof(uint8_t) * 64);

    for (int i = 0; i < 64; i++) {
        p->data[i] = IP[i];
    }

    return p;
}

Perm *create_permutation(uint8_t size, ...) {
    Perm *p = malloc(sizeof(Perm));
    p->size = size;
    p->data = calloc(size, sizeof(uint8_t));
    
    va_list va;
    va_start(va, size);

    for (int i = 0; i < size; i++) {
        p->data[i] = va_arg(va, int);
    }

    // Ending argument list traversal
    va_end(va);

    return p;
}

void free_permutation(Perm **p) {
    (*p)->size = 0;
    free((*p)->data);
    free(*p);

    // Safety measure, so that we know that all freed pointers are NULL.
    *p = NULL;
}

Perm *compute_inverse_permutation(Perm *perm) {
    Perm *p = malloc(sizeof(Perm));
    p->size = perm->size;
    p->data = calloc(perm->size, sizeof(uint8_t));
    
    for (int i = 0; i < perm->size; i++) {
        p->data[perm->data[i]-1] = i+1;
    }

    return p;
}

void print_permutation(Perm *perm) {
    for (int i = 0; i < perm->size; i++) {
        printf("%d ", perm->data[i]);
    }
    printf("\n");
}