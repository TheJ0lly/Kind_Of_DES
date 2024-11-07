#include "../include/permutation.h"
#include "../include/strconv.h"
#include "../include/errors.h"

#include <stdio.h>
#include <stdbool.h>
#include <getopt.h>
#include <string.h>

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
            // In cases where we need a value for a flag and we do not pass any value,
            // the next flag will be taken as the value, thus we check if the value is a flag.
            if (strcmp(optarg, "-h") == 0 || strcmp(optarg, "-t") == 0 || strcmp(optarg, "-d") == 0 || strcmp(optarg, "-e") == 0) {
                print_log(ERROR, "option -k requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
            }

            key = create_string(optarg, 8);

            if (key == NULL) {
                print_log(ERROR, "failed to allocate memory for key\n");
                return FAILED_ALLOC;
            }

            break;
        case 't':
            // We check if optarg is another flag, in which case we return error.
            // In cases where we need a value for a flag and we do not pass any value,
            // the next flag will be taken as the value, thus we check if the value is a flag.
            if (strcmp(optarg, "-h") == 0 || strcmp(optarg, "-k") == 0 || strcmp(optarg, "-d") == 0 || strcmp(optarg, "-e") == 0) {
                print_log(ERROR, "option -t requires an argument.\n");
                print_flags_description();
                return MISSING_ARG;
            }

            text = create_string(optarg, 256);

            if (text == NULL) {
                print_log(ERROR, "failed to allocate memory for text\n");
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
                print_log(ERROR, "option -%c requires an argument.\n", optopt);
                print_flags_description();
                return MISSING_ARG;
            }
            else {
                print_log(ERROR, "unknown option: %c\n", optopt);
                print_flags_description();
                return UNKNOWN_FLAG;
            }
        default:
            return -1;
        }
    }

    if (key == NULL) {
        print_log(ERROR, "key not passed - must pass an 8 character string.\n");
        return ARG_NOT_PASSED;
    }

    if (text == NULL) {
        print_log(ERROR, "text not passed - must pass an 8 character string.\n");
        return ARG_NOT_PASSED;
    }

    if (key->size != 8) {
        print_log(ERROR, "key must be exactly 8 characters.\n");
        return KEY_TOO_BIG;
    }

    // We check if the length of the text is divisible by 8.
    int textlen = check_len_by_8(text); 

    // If its not divisible by 8, we add padding.
    if (textlen != 0) {
        print_log(INFO, "text length is %lu - adding %d bytes of padding\n", text->size, textlen);
        if (add_padding(&text, textlen) != 0) {
            print_log(ERROR, "failed to allocate memory for string padding\n");
            return FAILED_ALLOC;
        }
        print_log(INFO, "padding successful - new length is %ld\n", text->size);
    }
    
    if (decrypt && encrypt) {
        print_log(ERROR, "both the encryption and decryption flags have been toggled - must have only 1.\n");
        return BOTH_OPS_TOGGLED;
    }

    if (!decrypt && !encrypt) {
        print_log(ERROR, "neither encryption nor decryption flag has been toggled - must have 1.\n");
        return BOTH_OPS_TOGGLED;
    }

}