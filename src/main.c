#include "../include/permutation.h"

#include <stdio.h>
#include <stdbool.h>
#include <getopt.h>
#include <string.h>

#define HELP                0
#define MISSING_ARG         1
#define UNKNOWN_FLAG        2
#define KEY_TOO_BIG         3
#define TEXT_TOO_BIG        4
#define BOTH_OPS_TOGGLED    5
#define NO_OPS_TOGGLED      6
#define ARG_NOT_PASSED      7

void print_flags_description() {
    printf("  -k <string>\n    the key to use for encryption/decryption.\n");
    printf("  -t <string>\n    the text to encrypt/decrypt.\n");
    printf("  -d\n    decrypt the text.\n");
    printf("  -e\n    encrypt the text.\n");
    printf("  -h\n    print the help menu.\n");
}

int main(int argc, char **argv) {
    char *key = NULL;
    char *text = NULL;
    bool encrypt = false;
    bool decrypt = false;

    char opt = 0;

    while ((opt = getopt(argc, argv, "hk:t:de")) != -1) {
        switch (opt) {
        case 'h':
            print_flags_description();
            return HELP;
        case 'k':
            key = optarg;

            // We check if optarg is another flag, in which case we return error.
            if (strcmp(key, "-h") == 0 || strcmp(key, "-t") == 0 || strcmp(key, "-d") == 0 || strcmp(key, "-e") == 0) {
                printf("option -k requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
            }
            break;
        case 't':
            text = optarg;

            // We check if optarg is another flag, in which case we return error.
            if (strcmp(key, "-h") == 0 || strcmp(key, "-k") == 0 || strcmp(key, "-d") == 0 || strcmp(key, "-e") == 0) {
                printf("option -t requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
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

    if (strlen(key) > 8) {
        printf("key is too big - max 8 characters allowed.\n");
        return KEY_TOO_BIG;
    }
    
    if (strlen(text) > 8) {
        printf("text is too big - max 8 characters allowed.\n");
        return TEXT_TOO_BIG;
    }

    if (decrypt && encrypt) {
        printf("both the encryption and decryption flags have been toggled - must have only 1.\n");
        return BOTH_OPS_TOGGLED;
    }

    if (!decrypt && !encrypt) {
        printf("neither encryption nor decryption flag has been toggled - must have 1.\n");
        return BOTH_OPS_TOGGLED;
    }

}