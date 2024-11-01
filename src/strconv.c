#include "../include/strconv.h"

#include <string.h>
#include <malloc.h>

string_t *create_string(char *str, int max) {
    string_t *nstr = NULL;
    
    if (str == NULL) {
        nstr = malloc(sizeof(string_t));
        
        if (nstr == NULL) {
            return NULL;
        }

        nstr->str = calloc(sizeof(char) * max, 1);

        if (nstr->str == NULL) {
            return NULL;
        }

        nstr->size = max;
        return nstr;
    }

    int index = 0;

    while (1) {
        if (index == max) {
            break;
        }

        if (str[index] == 0) {
            break;
        }

        index++;
    }

    nstr = malloc(sizeof(string_t));
    
    if (nstr == NULL) {
        return NULL;
    }

    nstr->str = malloc(sizeof(char) * index);

    if (nstr->str == NULL) {
        return NULL;
    }

    nstr->size = index;

    for (int i = 0; i < index; i++) {
        nstr->str[i] = str[i];
    }

    return nstr;
}

void free_string(string_t **str) {
    free((*str)->str);
    free((*str));

    // To be safe.
    *str = NULL;
}

void print_string(string_t *str, char *format) {
    int fmtlen = strlen(format);

    for (int i = 0; i < fmtlen; i++) {
        if (format[i] == '%' && format[i+1] == 's') {
            for (int i = 0; i < str->size; i++) {
                printf("%c", str->str[i]);
            }
            // We increment i to get rid of s.
            i++;
        } else {
            printf("%c", format[i]);
        }
    }
}

int check_len_by_8(string_t *str) {
    int target = 8;

    while (str->size > target) {
        target += 8;
    }

    return target - str->size;
}

int add_padding(string_t **str, int pad) {
    int len = (*str)->size;
    char *temp = malloc(sizeof(char) * (len + 1));

    if (temp == NULL) {
        return 1;
    }

    for (int i = 0; i < len; i++) {
        temp[i] = (*str)->str[i];
    }

    // We free the initial string.
    free_string(str);

    // We reallocate the string on the heap.
    *str = create_string(NULL, len + pad);

    if (*str == NULL) {
        return 1;
    }

    for (int i = 0; i < len; i++) {
        (*str)->str[i] = temp[i];
    }

    free(temp);

    return 0;
}

uint64_t get_bitmap64(string_t *str, int start) {
    uint64_t toret = 0;

    for (int i = start; i < str->size; i++) {
        toret <<= 8; 
        toret |= str->str[i];
    }

    return toret;
}