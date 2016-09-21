#ifndef SUBHASH_H
#define SUBHASH_H 0
#include "list.h"

typedef struct __pair__ {
    char* key;
    char* value;
    struct __pair__* next;
} PAIR;

PAIR* new_pair()
{
    PAIR* pair = (PAIR*) malloc(sizeof(PAIR));

    pair->key   = NULL;
    pair->value = NULL;
    pair->next  = NULL;

    return pair;
}

PAIR* butt(PAIR* head)
{
    PAIR* pair = head;

    while (pair->next != NULL)
        inc(pair);

    return pair;
}

PAIR* associate(PAIR* head, char* key, char* value)
{
    PAIR* pair = butt(head);

    pair->next = new_pair();
    pair = pair->next;
    pair->key = key;
    pair->value = value;

    return head;
}

char* find_in_pair(PAIR* head, char* to_find)
{
    PAIR* pair = head->next;
    char* outlet = NULL;
    int index = 0;

    while ((pair != NULL) && (outlet == NULL))
    {
        if (compare(to_find, pair->key) == EQUAL)
            outlet = pair->value;

        inc(pair);
    }

    return outlet;
}

void write_pair(PAIR* head)
{
    PAIR* pair = head->next;

    if (pair != NULL) {
        printf("%s:%s", pair->key, pair->value);
        inc(pair);
    }
    while (pair != NULL)
    {
        printf(" %s:%s", pair->key, pair->value);
        inc(pair);
    }
}

char* delete_pair(PAIR* head, char* to_find)
{
    PAIR* list  = head;
    PAIR* del   = NULL;
    char* key   = NULL;
    char* value = NULL;

    while (list->next != NULL)
    {
        key = (list->next)->key;

        if (compare(to_find, key) == EQUAL)
            break;

        inc(list);
    }

    if (list->next != NULL) {
        del = list->next;
        list->next = del->next;
        value = del->value;
        free(del);
    }

    return value;
}

#endif
