#include "wc.h"

int parseFlags(char *input) {
  uint8_t flags = 0;

  while (*input) {
    switch (*input) {
    case 'c':
      flags |= ByteFlag;
      flags &= ~MultFlag;
      break;
    case 'l':
      flags |= LineFlag;
      break;
    case 'w':
      flags |= WordFlag;
      break;
    case 'm':
      flags |= MultFlag;
      flags &= ~ByteFlag;
      break;
    }

    input++;
  }

  return flags;
}

int isFlagSet(flags_t flag, int flags) { return (flags & flag) != 0; }

void printer(counter_t *count, char *file, int flags) {
  if (isFlagSet(LineFlag, flags)) {
    printf("%7d\t", count->line);
  }

  if (isFlagSet(WordFlag, flags)) {
    printf("%7d\t", count->word);
  }

  if (isFlagSet(ByteFlag, flags)) {
    printf("%7d\t", count->byte);
  }

  if (isFlagSet(MultFlag, flags)) {
    printf("%7d\t", count->multibyte);
  }

  printf("%s\n", file);
}

void counter(char *buf, unsigned int length, counter_t *count) {
  for (int i = 0; i < length; i++) {
    char ch = buf[i];

    count->byte++;

    if (ch == '\n') {
      count->line++;
    }

    if ((ch & 0xC0) != 0x80) {
      count->multibyte++;
    }

    if (isspace(ch)) {
      count->in_word = 0;
    } else if (!count->in_word) {
      count->in_word = 1;
      count->word++;
    }
  }
}

void countFromFile(char *file, int flags, counter_t *count) {
  int fd = open(file, O_RDONLY);
  if (fd < 0) {
    close(fd);

    perror("file errors");
    return;
  }

  int size;
  char buf[1024 * 4];

  while (1) {
    size = read(fd, buf, sizeof(buf));
    if (size > 0)
      counter(buf, size, count);
    else if (size == 0)
      break;
    else {
      perror("read error");
      break;
    }
  }

  close(fd);
}

int main(int argc, char *argv[]) {
  if (argc == 1) {
    printf("Usage: cwc -[clw] [file ...]");
    return 1;
  }

  int flags = WordFlag | LineFlag | ByteFlag;
  int offset = 1;

  if (*argv[1] == '-') {
    offset = 2;
    flags = parseFlags(argv[1]);
  }

  // if (argc <= 2) {
  //   countFromStdin();
  // }

  int files = argc - offset;
  counter_t sum = {0};

  for (int i = offset; i < argc; i++) {
    counter_t count = {0};

    countFromFile(argv[i], flags, &count);
    printer(&count, argv[i], flags);

    sum.line += count.line;
    sum.byte += count.byte;
    sum.word += count.word;
    sum.multibyte += count.multibyte;
  }

  if (files != 1) {
    printer(&sum, "total", flags);
  }

  return 0;
}
