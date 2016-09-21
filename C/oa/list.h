#ifndef LIST_H
#define LIST_H

#define inc(A) ((A) = (A)->next)

typedef struct NODE
{
    struct NODE* next;
    char* key;
    char* value;
} NODE;

typedef NODE LIST;

LIST* list_new()
{
    LIST* list = (LIST*) malloc(sizeof(LIST));

    list->next = NULL;
    list->key = NULL;
    list->value = NULL;

    return list;
}
LIST* new_list() { return list_new(); }

LIST* tail(LIST* head)
{
    LIST* list = head;

    while (list->next != NULL)
        list = list->next;

    return list;
}

int list_length(LIST* head)
{
    LIST* list   = head->next;
    int   result = 0;

    while (list != NULL)
    {
        inc(list);
        result++;
    }

    return result;
}

char* list_get(LIST* head, int index)
{
    LIST* list = head->next;
    char* outlet = NULL;
    int i = 0;

    while ((i != index) && (list != NULL))
    {
        inc(list);
        ++i;
    }

    if (list != NULL)
        outlet = list->value;

    return outlet;
}

LIST* list_add(LIST* head, char* data)
{
    LIST* new_node = new_list();
    LIST* list = tail(head);

    new_node->value = data;
    list->next = new_node;

    return head;
}

LIST* list_add_at(LIST* head, char* data, int index)
{
    LIST* new_node = new_list();
    LIST* list = head;
    int i = 0;

    new_node->value = data;
    while (list != NULL)
    {
        if (i == index) {
            new_node->next = list->next;
            list->next = new_node->next;
            break;
        }

        inc(list);
        i++;
    }

    if (list == NULL) {
        list = tail(head);
        list->next = new_node;
        new_node->next = NULL;
    }

    return head;
}

LIST* associate(LIST* head, char* key, char* value)
{
    LIST* pair = tail(head);

    pair->next = new_list();
    pair = pair->next;
    pair->key = key;
    pair->value = value;

    return head;
}

int list_find(LIST* head, char* to_find)
{
    LIST* list = head->next;
    int outlet = -1;
    int index = 0;

    while ((list != NULL) && (outlet < 0))
    {
        if (equals(to_find, list->value))
            outlet = index;

        ++index;
        inc(list);
    }

    return outlet;
}
char* pair_find(LIST* head, char* to_find)
{
    LIST* pair = head->next;
    char* outlet = NULL;

    while ((pair != NULL) && (outlet == NULL))
    {
        if (compare(to_find, pair->key) == EQUAL)
            outlet = pair->value;

        inc(pair);
    }

    return outlet;
}

/*
char* to_string(LIST* head)
{
    LIST* list = head->next;
    char* outlet = "";

    outlet = concat(outlet, "---\n- list:\n");
    while (list != NULL)
    {
        outlet = concat(outlet, "  - ");
        outlet = concat(outlet, list->value);
        outlet = concat(outlet, "\n");
        inc(list);
    }

    outlet = concat(outlet, "...\n");
    return outlet;
}
*/
char* list_to_string(LIST* head)
{
    LIST* list = head->next;
    char* outlet = "";

    while (list != NULL)
    {
        outlet = concat(outlet, list->value);
        outlet = concat(outlet, "\n");
        inc(list);
    }

    return outlet;
}

char* list_yaml_with_title(LIST* head, char* title)
{
    LIST* list = head->next;
    char* outlet = "";

    outlet = concat(outlet, "---\n- ");
    outlet = concat(outlet, title);
    outlet = concat(outlet, ":\n");

    while (list != NULL)
    {
        outlet = concat(outlet, "  - ");
        outlet = concat(outlet, list->value);
        outlet = concat(outlet, "\n");
        inc(list);
    }

    outlet = concat(outlet, "...\n");
    return outlet;
}

char* list_yaml(LIST* head)
{
    LIST* list = head->next;
    char* yaml = "---\n";

    while (list != NULL)
    {
        cat(yaml, "- ");
        if (list->value != NULL)
            yaml = concat(yaml, list->value);
        else
            yaml = concat(yaml, "NULL");
        cat(yaml, "\n");
        inc(list);
    }

    cat(yaml, "...\n");
    return yaml;
}

void write_list_to_file(LIST* head, char* output)
{
    FILE* outlet = fopen(output, "w");
    LIST* list   = head->next;

    for (list = head->next; list != NULL; inc(list))
        fprintf(outlet, "%s\n", list->value);

    fclose(outlet);
}

char* pair_to_string(LIST* head)
{
    LIST* pair = head->next;
    char* result = "";

    while (pair != NULL)
    {
        cat(result, "- ");
        cat(result, pair->key);
        cat(result, ":");
        cat(result, pair->value);
        cat(result, "\n");
        inc(pair);
    }

    return result;
}

void list_print(LIST* head)
{
    LIST* list   = head->next;

    for (list = head->next; list != NULL; inc(list))
        printf("%s\n", list->value);
}

int list_contains(LIST* head, char* to_find)
{
    LIST* list = head->next;
    int result = 0;

    while ((list != NULL) && (result == 0))
    {
        result = equals(to_find, list->value);
        inc(list);
    }

    return result;
}

LIST* using_bubblesort(LIST* head, int(*function)(int))
{
    LIST* list = head->next;
    char* temp = NULL;
    char* stri = NULL;
    char* strj = NULL;
    int   len  = list_length(head);
    int   i    = 0;
    int   flag = 1;

    while (flag)
    {
        flag = 0;

        for (i = 0, list = head->next; i < len - 1; ++i, inc(list))
        {
            stri = list->value;
            strj = (list->next)->value;

            if (function(compare(stri, strj)))
            {
                temp = (list->next)->value;
                (list->next)->value = list->value;
                list->value = temp;
                flag = 1;
            }
        }
    }

    /* MEOW */
    if ((head->next)->value == NULL)
        head->next = (head->next)->next;

    return head;
}

int is_bigger(int result)
{
    if (result == BIGGER)
        return 1;
    else
        return 0;
}

LIST* list_sort(LIST* head)
{
    return using_bubblesort(head, &is_bigger);
}

char* list_pop(LIST* head)
{
    LIST* list   = head->next;
    char* result = NULL;

    if (list != NULL)
    {
        head->next = list->next;
        result = list->value;
    }

    return result;
}

LIST* list_push(LIST* head, char* item)
{
    LIST* list = tail(head);
    list->next = (LIST*) malloc(sizeof(LIST));
    inc(list);
    list->value = item;
    list->next = NULL;
    return head;
}
#define push(A,B) ((A) = list_push((A),(B)))

LIST* list_remove(LIST* head, int index)
{
    LIST* list   = head->next;
    LIST* to_del = NULL;
    int i = 0;

    while ((i + 1 != index) && (list != NULL))
    {
        inc(list);
        ++i;
    }

    if (list->next != NULL)
    {
        to_del = list->next;
        list->next = to_del->next;
        free(to_del);
    }

    return head;
}

char* pair_delete(LIST* head, char* to_find)
{
    LIST* list  = head;
    LIST* del   = NULL;
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

void list_free(LIST* head)
{
    LIST* list = head;
    LIST* memo;

    while (list->next != NULL)
    {
        memo = list->next;
        free(list);
        list = memo;
    }

    free(list);
}

LIST* list_filter(LIST *list, bool (*func)(char*))
{
    LIST *filtered = list_new();
    LIST *it = NULL;

    for (it = list->next; it != NULL; ++it)
        if ((*func)(it->value) == true)
            push(filtered, it->value);

    return filtered;
}

LIST* list_strsplit(char* string, char to_divide)
{
    LIST* result  = new_list();
    char* section = "";
    int   length  = strlen(string);
    int   i       = 0;

    for (i = 0; i < length; ++i)
    {
        section = "";

        while (string[i] != to_divide && i < length)
        {
            section = concat(section, to_array(string[i]));
            i++;
        }

        if (strlen(section) > 0)
            push(result, section);
    }

    return result;
}

#endif
