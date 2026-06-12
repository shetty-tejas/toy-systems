#ifndef CWC
#define CWC

#include <wctype.h>
#include <fcntl.h>
#include <stdint.h>
#include <stdio.h>
#include <unistd.h>

typedef enum flags {
  ByteFlag = 0b0001,
  LineFlag = 0b0010,
  WordFlag = 0b0100,
  MultFlag = 0b1000,
} flags_t;

typedef struct counter {
  int byte;
  int line;
  int word;
  int multibyte;

  int in_word;
} counter_t;

#endif
