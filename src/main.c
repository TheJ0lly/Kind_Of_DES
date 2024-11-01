#include "../include/permutation.h"
#include "../include/strconv.h"


#include <stdio.h>
#include <stdbool.h>
#include <getopt.h>
#include <string.h>
#include <inttypes.h>

#define HELP                0
#define MISSING_ARG         1
#define UNKNOWN_FLAG        2
#define KEY_TOO_BIG         3
#define TEXT_TOO_BIG        4
#define BOTH_OPS_TOGGLED    5
#define NO_OPS_TOGGLED      6
#define ARG_NOT_PASSED      7
#define FAILED_ALLOC        8

void print_flags_description() {
    printf("  -k <string>\n    the key to use for encryption/decryption.\n");
    printf("  -t <string>\n    the text to encrypt/decrypt.\n");
    printf("  -d\n    decrypt the text.\n");
    printf("  -e\n    encrypt the text.\n");
    printf("  -h\n    print the help menu.\n");
}

int main(int argc, char **argv) {
    string_t *key = NULL;
    string_t *text = NULL;
    bool encrypt = false;
    bool decrypt = false;


    char opt = 0;

    while ((opt = getopt(argc, argv, "hk:t:de")) != -1) {
        switch (opt) {
        case 'h':
            print_flags_description();
            return HELP;
        case 'k':
            // We check if optarg is another flag, in which case we return error.
            if (strcmp(optarg, "-h") == 0 || strcmp(optarg, "-t") == 0 || strcmp(optarg, "-d") == 0 || strcmp(optarg, "-e") == 0) {
                printf("option -k requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
            }

            key = create_string(optarg, 8);

            if (key == NULL) {
                printf("error: failed to allocate memory for key\n");
                return FAILED_ALLOC;
            }

            break;
        case 't':
            // We check if optarg is another flag, in which case we return error.
            if (strcmp(optarg, "-h") == 0 || strcmp(optarg, "-k") == 0 || strcmp(optarg, "-d") == 0 || strcmp(optarg, "-e") == 0) {
                printf("option -t requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
            }

            text = create_string(optarg, 256);

            if (text == NULL) {
                printf("error: failed to allocate memory for text\n");
                return FAILED_ALLOC;
            }

            break;
        case 'e':
            encrypt = true;
            break;
        case 'd':
            decrypt = true;
            break;
        case '?':
            if (optopt == 'k' || optopt == 't') {
                printf("option -%c requires an argument.\n", optopt);
                print_flags_description();
                return MISSING_ARG;
            }
            else {
                printf("unknown option: %c\n", optopt);
                print_flags_description();
                return UNKNOWN_FLAG;
            }
        default:
            return -1;
        }
    }

    if (key == NULL) {
        printf("key not passed - must pass an 8 character string.\n");
        return ARG_NOT_PASSED;
    }

    if (text == NULL) {
        printf("text not passed - must pass an 8 character string.\n");
        return ARG_NOT_PASSED;
    }

    if (key->size != 8) {
        printf("key must be exactly 8 characters.\n");
        return KEY_TOO_BIG;
    }

    int textlen = check_len_by_8(text); 
    if (textlen != 0) {
        printf("text length is %ld - adding %d bytes of padding\n", text->size, textlen);
        if (add_padding(&text, textlen) != 0) {
            printf("error: failed to allocate memory for string padding\n");
            return FAILED_ALLOC;
        }
        printf("padding successful - new length is %ld\n", text->size);
        print_string(text, "string: %s\n");
    }

    printf("bitmap for text: %lu\n", get_bitmap64(text, 0));
    
    if (decrypt && encrypt) {
        printf("both the encryption and decryption flags have been toggled - must have only 1.\n");
        return BOTH_OPS_TOGGLED;
    }

    if (!decrypt && !encrypt) {
        printf("neither encryption nor decryption flag has been toggled - must have 1.\n");
        return BOTH_OPS_TOGGLED;
    }

}