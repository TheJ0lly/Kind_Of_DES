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

    va_list va;
    va_start(va, format);
    vprintf(format, va);
    va_end(va);
}