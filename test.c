#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <malloc.h>
#include <threads.h>
#include <strings.h>

#define SIZE 10
#define MEOW 10

void increment(int *p);

int main(void)
{
  printf("an int uses %zu bytes of memory\n", sizeof(int));

  int i = 10;
  int *p;

  p = &i;

  printf("value of i is %d\n", i);
  printf("address if i is %p\n", (void *)p);

  *p = 20;

  printf("i is %d\n", i);
  printf("i is %d\n", *p);

  i = 10;
  int *j = &i;
  printf("i is %d\n", i);
  printf("i is also %d\n", *j);

  increment(j);

  printf("i is %d\n", i);

  int *np;
  np = NULL;
}

void increment(int *p)
{
  *p = *p + 1;
}
