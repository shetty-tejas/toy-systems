#include <ctype.h>
#include <fcntl.h>
#include <stdint.h>
#include <stdio.h>
#include <unistd.h>

const uint8_t BYTE_FLAG = 0b001;
const uint8_t LINE_FLAG = 0b010;
const uint8_t WORD_FLAG = 0b100;

#define BYTE_FLAG_SET(flags) (((flags) & BYTE_FLAG) != 0)
#define LINE_FLAG_SET(flags) (((flags) & LINE_FLAG) != 0)
#define WORD_FLAG_SET(flags) (((flags) & WORD_FLAG) != 0)

int parseFlags(char *input) {
  uint8_t flags = 0;

  for (int i = 0; input[i] != '\0'; i++) {
    int set = 0;

    switch (input[i]) {
    case 'c':
      set = BYTE_FLAG;
      break;
    case 'l':
      set = LINE_FLAG;
      break;
    case 'w':
      set = WORD_FLAG;
      break;
    }

    flags |= set;
  }

  return flags;
}

void printer(int lineCount, int wordCount, int byteCount, char *file,
             int flags) {
  if (LINE_FLAG_SET(flags)) {
    printf("%7d\t", lineCount);
  }

  if (WORD_FLAG_SET(flags)) {
    printf("%7d\t", wordCount);
  }

  if (BYTE_FLAG_SET(flags)) {
    printf("%7d\t", byteCount);
  }

  printf("%s\n", file);
}

void counter(char buf, int *byteCount, int *lineCount, int *inWord,
             int *wordCount) {
  (*byteCount)++;

  if (buf == '\n') {
    (*lineCount)++;
  }

  if (isspace(buf)) {
    (*inWord) = 0;
  } else if (!*inWord) {
    (*inWord) = 1;
    (*wordCount)++;
  }
}

void countFromFile(char *file, int flags) {
  int fd = open(file, O_RDONLY);
  if (fd < 0) {
    perror("file errors");
    return;
  }

  int byteCount = 0, lineCount = 0, wordCount = 0, inWord = 0;
  char buf[1024 * 4];

  while (1) {
    int rs = read(fd, buf, sizeof(buf));
    if (rs <= 0) {
      if (rs < 0) {
        perror("read error");
      }

      break;
    }

    for (int i = 0; i < rs; i++) {
      counter(buf[i], &byteCount, &lineCount, &inWord, &wordCount);
    }
  }

  printer(lineCount, wordCount, byteCount, file, flags);
  close(fd);
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    printf("Usage: cwc -[c] [file ...]");
    return 1;
  }

  const int flags = parseFlags(argv[1]);

  if (argc >= 3) {
    for (int i = 2; i < argc; i++) {
      countFromFile(argv[i], flags);
    }
  } else {
    printf("not implemented");
  }
}
