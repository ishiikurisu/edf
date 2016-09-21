#pragma once
#include <stdlib.h>
#include <stdio.h>
#include <oa.h>

typedef struct
{
	char *id;
	FILE *outlet;
	int position;

} CHANNEL;


DICT* process_labels(char* raw_data)
{
    LIST* labels = list_strsplit(raw_data, ' ');
    LIST* label = NULL;
    DICT* enumeration = new_map(7, dumb);
    char* indexstr;
    int index = 0;

    for (label = labels->next; label != NULL; inc(label), ++index)
    {
        indexstr = itos(index);
        map_put(enumeration, label->value, indexstr);
    }

    return enumeration;
}
DICT* process_column(char* column)
{
    DICT* labels = NULL;
    LIST* data = list_strsplit(column, ':');
    char* key = list_get(data, 0);
    char* value = list_get(data, 1);

    if (compare(key, " labels") == EQUAL) {
        labels = process_labels(value);
    }

    return labels;
}
DICT* get_labels(FILE* csv)
{
    char* header = read_from_file(csv);
    LIST* columns = list_strsplit(header, ',');
    LIST* column = NULL;
    DICT* labels = NULL;

    for (column = columns->next; column != NULL && labels == NULL; inc(column))
    {
        labels = process_column(column->value);
    }

    return labels;
}

DICT* get_needed_labels(DICT *all, char **channels)
{
    DICT *labels = new_map(3, dumb);
    char *key, *value;
    int i = 0;

    for (i = 0; i < 4; ++i)
    {
        key = channels[i];
        value = map_get(all, key);
        labels = map_put(labels, key, value);
    }

    return labels;
}
