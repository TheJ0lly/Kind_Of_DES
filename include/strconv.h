#ifndef _STRCONV_H_
#define _STRCONV_H_

#include <stdint.h>
#include <stdio.h>



typedef struct string_t {
    char *str;
    size_t size;
} string_t;

// We create a new string_t from a C-string.
// If `str` is NULL, we simply allocate a string of `max` size.
// If the length cannot be determined, `max` will be used as an enforced limit.
string_t *create_string(char *str, int max);

// Frees the string_t.
void free_string(string_t **str);

// Prints the string_t.
void print_string(string_t *str, char *format);

// Get the string as a C-string, with the NULL-terminator character at the end.
// Acts like sprintf.
char *get_Cstring(string_t *str);

// Returns the difference, if any, that is needed until the next bigger multiple of 8.
int check_len_by_8(string_t *str);

// Adds padding to the string, so that the length is divisible by 8.
int add_padding(string_t **str, int pad);

// Returns the bitmap of 8 bytes from `start`.
uint64_t get_bitmap64(string_t *str, int start);

#endif