#include <stdio.h>
#include <malloc.h>
#include <string.h>


int main() {
    char *h = malloc(sizeof(char) * strlen("Hello, World"));
    h = "Hello, World";

    printf("%s\n", h);
}