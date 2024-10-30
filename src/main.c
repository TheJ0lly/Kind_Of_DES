#include "../include/permutation.h"

int main(int argc, char **argv) {
    Perm *ip = get_default_initial_permutation();

    Perm *inv_ip = compute_inverse_permutation(ip);

    print_permutation(ip);
    print_permutation(inv_ip);

    free_permutation(&ip);
    free_permutation(&inv_ip);
}