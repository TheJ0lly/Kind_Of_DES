#include "../include/errors.h"

#include <stdarg.h>
#include <stdio.h>
#include <string.h>

void print_log(ERR_LEVEL level, char *format, ...) {
    switch (level) {
    case INFO:
        printf("info: ");
        break;
    case WARNING:
        printf("warning: ");
        break;
    case ERROR:
        printf("error: ");
        break;
    default:
        printf("unknown log level: ");
        break;
    }

    int fmtlen = strlen(format);

    va_list va;
    va_start(va, format);

    for (int i = 0; i < fmtlen; i++) {
        if (format[i] == '%') {
            switch (format[i+1]) {
            case 's':
                printf("%s", va_arg(va, char*));
                break;
            case 'd':
                printf("%d", va_arg(va, int));
                break;
            case 'c':
                printf("%c", va_arg(va, int));
                break;
            case 'l':
                switch (format[i+2])
                {
                case 'd':
                    printf("%ld", va_arg(va, long));
                    break;
                case 'u':
                    printf("%lu", va_arg(va, unsigned long));
                    break;
                default:
                    printf(" --- UNKNOWN FORMAT: %c ---", format[i+2]);
                    break;
                }

                // We increment to get rid of the modifying character.
                i++;
                break;
            default:
                printf(" --- UNKNOWN FORMAT: %c ---", format[i+1]);
                break;
            }
            // We increment to get rid of the modifying character.
            i++;
        } else {
            printf("%c", format[i]);
        }

    }
}