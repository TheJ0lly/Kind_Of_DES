#ifndef _ERRORS_H_
#define _ERRORS_H_

typedef enum ERROR_LEVEL {
    INFO = 0,
    WARNING,
    ERROR,
} ERR_LEVEL;

void print_log(ERR_LEVEL level, char *format, ...);

#endif