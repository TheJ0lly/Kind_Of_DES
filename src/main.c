#include <stdio.h>
#include <malloc.h>
#include <string.h>


int main(int argc, char **argv) {
    char *h = malloc(sizeof(char) * strlen("Hello, World"));
    h = "Hello, World";

    printf("%s\n", h);

    free(h);
}