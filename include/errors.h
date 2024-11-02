#ifndef _ERRORS_H_
#define _ERRORS_H_

#define HELP                0
#define MISSING_ARG         1
#define UNKNOWN_FLAG        2
#define KEY_TOO_BIG         3
#define TEXT_TOO_BIG        4
#define BOTH_OPS_TOGGLED    5
#define NO_OPS_TOGGLED      6
#define ARG_NOT_PASSED      7
#define FAILED_ALLOC        8

typedef enum ERROR_LEVEL {
    INFO = 0,
    WARNING,
    ERROR,
} ERR_LEVEL;

void print_log(ERR_LEVEL level, char *format, ...);

#endif