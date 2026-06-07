#include <stdio.h>
#include <stdlib.h>

#define linkedlist struct linked_list

struct node {
  int value;

  struct node *prev;
  struct node *next;
};

struct node *newnode(int val, struct node *prev, struct node *next) {
  struct node *ptr = malloc(sizeof(struct node));

  if (ptr != NULL) {
    *ptr = ((struct node){.value = val, .prev = prev, .next = next});
  }

  return ptr;
}

void destroynode(struct node *node) { free(node); }

struct linked_list {
  struct node *head;
  struct node *tail;

  int size;
};

struct linked_list newlist() {
  return (linkedlist){.size = 0, .head = NULL, .tail = NULL};
}

void push(linkedlist *ll, int val) {
  if (ll->tail == NULL) {
    ll->head = ll->tail = newnode(val, NULL, NULL);
  } else {
    ll->tail->next = newnode(val, ll->tail, NULL);
    ll->tail = ll->tail->next;
  }

  ll->size++;
}

void pop(linkedlist *ll) {
  struct node *node = ll->tail;
  if (node == NULL) {
    return;
  }

  if (ll->tail == ll->head) {
    ll->head = ll->tail = NULL;
  } else {
    ll->tail = ll->tail->prev;
    ll->tail->next = NULL;
  }

  ll->size--;
  destroynode(node);
}

void shift(linkedlist *ll, int val) {
  if (ll->head == NULL) {
    ll->head = ll->tail = newnode(val, NULL, NULL);
  } else {
    ll->head = newnode(val, NULL, ll->head);
    ll->head->next->prev = ll->head;
  }

  ll->size++;
}

void unshift(linkedlist *ll) {
  struct node *node = ll->head;
  if (node == NULL) {
    return;
  }

  if (ll->head == ll->tail) {
    ll->head = ll->tail = NULL;
  } else {
    ll->head = ll->head->next;
    ll->head->prev = NULL;
  }

  ll->size--;
  destroynode(node);
}

int insert(linkedlist *ll, int val, int index) {
  if (index <= 0) {
    shift(ll, val);
    return 0;
  } else if (index >= ll->size) {
    push(ll, val);
    return ll->size - 1;
  }

  struct node *curr = ll->head;
  int x = index;

  while (x > 0) {
    curr = curr->next;
    x--;
  }

  struct node *prev = curr->prev;
  struct node *new = newnode(val, prev, curr);

  prev->next = curr->prev = new;

  ll->size++;
  return index;
}

int delete (linkedlist *ll, int index) {
  if (index <= 0) {
    unshift(ll);
    return 0;
  } else if (index >= ll->size - 1) {
    pop(ll);
    return ll->size - 1;
  }

  struct node *curr = ll->head;
  int x = index;

  while (x > 0) {
    curr = curr->next;
    x--;
  }

  struct node *prev = curr->prev;
  struct node *next = curr->next;

  if (prev != NULL) {
    prev->next = next;
  }

  if (next != NULL) {
    next->prev = prev;
  }

  destroynode(curr);
  ll->size--;

  return index;
}

void printll(linkedlist *ll) {
  printf("---printing ll---\n");

  struct node *curr = ll->head;
  int pos = 0;

  while (curr != NULL) {
    printf("%d: %d\n", pos, curr->value);

    pos++;
    curr = curr->next;
  }

  printf("---printing done---\n");
}

int main(void) {
  linkedlist ll = newlist();

  push(&ll, 2);
  push(&ll, 4);
  push(&ll, 6);
  push(&ll, 8);
  push(&ll, 10);

  insert(&ll, 0, 0);
  printll(&ll);
  insert(&ll, 1, 1);
  printll(&ll);
  insert(&ll, 3, 3);
  printll(&ll);
  insert(&ll, 5, 5);
  printll(&ll);
  insert(&ll, 7, 7);
  printll(&ll);
  insert(&ll, 9, 9);
  printll(&ll);

  delete (&ll, 9);
  printll(&ll);
  delete (&ll, 7);
  printll(&ll);
  delete (&ll, 5);
  printll(&ll);
  delete (&ll, 3);
  printll(&ll);
  delete (&ll, 1);
  printll(&ll);
  delete (&ll, 0);
  printll(&ll);

  insert(&ll, 12, ll.size);
  printll(&ll);
  delete (&ll, ll.size);
  printll(&ll);
  delete (&ll, ll.size - 1);
  printll(&ll);
}
