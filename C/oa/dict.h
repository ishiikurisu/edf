#ifndef HASH_H
#define HASH_H 0

/* please include the list.h header */

typedef struct {
    int size;
    LIST** data;
    int (*func)(char*);
    int (*hash)(char*);
} DICT;

DICT* new_map(int size, int (*hash)(char*))
{
    DICT* map = (DICT*) malloc(sizeof(DICT));
    LIST** data = (LIST**) malloc(size * sizeof(LIST*));
    int i = 0;

    for (i = 0; i < size; ++i)
        data[i] = new_list();

    map->size = size;
    map->data = data;
    map->func = hash;
    map->hash = hash;

    return map;
}
DICT* map_new(int sz, int(*h)(char*)) { return new_map(sz,h); }

int get_address(DICT* map, char* key)
{
    return map->hash(key) % map->size;
}

DICT* put_in_map(DICT* map, char* key, char* value)
{
    if (map == NULL || key == NULL || value == NULL) return map;
    int index = get_address(map, key);
    map->data[index] = associate(map->data[index], key, value);
    return map;
}
DICT* map_put(DICT* m, char* k, char* v){ return put_in_map(m, k, v); }
DICT* map_add(DICT* m, char* k, char* v){ return put_in_map(m, k, v); }

char* get_from_map(DICT* map, char* key)
{
    if (map == NULL || key == NULL) return NULL;
    return pair_find(map->data[get_address(map, key)], key);
}
char* map_get(DICT* m, char* k) { return get_from_map(m, k); }

int map_contains_key(DICT* map, char* key)
{
    if (map == NULL || key == NULL)
        return 0;

    if (get_from_map(map, key) != NULL)
        return 1;
    else
        return 0;
}
int map_contains(DICT* m, char* k) { return map_contains_key(m, k); }

void map_print(DICT* map)
{
    char* str = NULL;
    int i = 0;

    for (i = 0; i < map->size; ++i)
    {
        str = pair_to_string(map->data[i]);

        printf("%d: \n", i);
        printf("%s", str);
        printf("\n");
    }
}
char* map_to_string(DICT* map)
{
    char* result = "---\n";
    char* pairs  = NULL;
    int i = 0;

    for (i = 0; i < map->size; ++i)
    {
        pairs = pair_to_string(map->data[i]);

        cat(result, itos(i));
        cat(result, ": \n");
        cat(result, pairs);
    }

    cat(result, "...\n");
    return result;
}

char* map_yaml(DICT* m) { return map_to_string(m); }

DICT* map_feed(DICT* map, char* input, char sep)
{
    FILE* inlet  = fopen(input, "r");
    char* line = read_from_file(inlet);
    LIST* data = NULL;
    char* key  = NULL;
    char* value = NULL;

    while (!feof(inlet))
    {
        data = list_strsplit(line, sep);
        key = list_get(data, 0);
        value = list_get(data, 1);
        map = put_in_map(map, key, value);
        line = read_from_file(inlet);
    }

    fclose(inlet);
    return map;
}

char* remove_key(DICT* map, char* key)
{
    int index = get_address(map, key);
    return pair_delete(map->data[index], key);
}
char* map_remove(DICT* m, char* k){ return remove_key(m, k); }


#endif
