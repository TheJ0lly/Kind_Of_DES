#ifndef _PERMUTATION_H_
#define _PERMUTATION_H_

#include <stdint.h>

typedef struct Permutation {
    int *data;
    int size;
} Perm;

// The hard-coded IP from permutation.c.
Perm *get_default_initial_permutation();

// Create and free permutations.
Perm *create_permutation(int size, ...);

// We take a pointer to a pointer to a Perm, to free the underlying array, and also the Perm pointer.
void free_permutation(Perm **p);


Perm *compute_inverse_permutation(Perm *perm);


// Helper function.
void print_permutation(Perm *perm);



#endif