#include <stdint.h>
#include <stdio.h>
#include <string.h>

void compress() {
  int16_t curr = getc(stdin);
  int16_t next = curr;
  uint8_t count = 0;

  do {
    if (next == curr && count < 0xFF) {
      count++;
    } else {
      if (count > 0) {
        putc(curr, stdout);
        putc(count, stdout);
      }

      curr = next;
      count = 1;
    }

    next = getc(stdin);
  } while (curr != EOF);
}

void decompress() {
  int16_t curr;
  int16_t count;

  while (1) {
    curr = getc(stdin);
    count = getc(stdin);

    if (curr == EOF || count == EOF) {
      break;
    }

    while (count-- > 0) {
      putc(curr, stdout);
    }
  };
}

int main(int argc, char *argv[]) {
  if (argc != 2) {
    printf("Usage: rle [compress|decompress]");
    return 1;
  }

  if (!strcmp(argv[1], "compress")) {
    compress();
  } else if (!strcmp(argv[1], "decompress")) {
    decompress();
  } else {
    printf("Usage: rle [compress|decompress]");
    return 1;
  }

  return 0;
}
