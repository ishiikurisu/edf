# String

``` C
int strlen(char*);
char* concat(char*, char*);
char* ctos(char);
char* itos(int);
char* substring(char*, int, int);
int compare(char*, char*);
char last_char(char*);
char* tidy_string(char*);
char** strsplit(char*,char); /* under construction */
```

# Linked data types

Every function in those headers use the following convention:

    type_method(obj, args...)

## Available LDTs

+ Lists
+ Maps
+ Sets -- under construction

## Methods

+ new
+ length
+ get
+ add
+ find
+ print
  - yaml
+ contains
+ to_string
+ remove
+ free

## Specific to every type

+ Lists
  - add_at
  - pop
  - push
  - sort
  - split(char*, char)
  - associate
  - strsplit

+ Maps
  - put

Cosequential list processing
============================

Can be added by the `coseq.h` header.

Methods
-------

``` C
char* read_from_file(FILE* fp);
char* get_line_from_file(FILE* inlet, int line_number);
char* get_line(char* input, int line_number);
LIST* read_whole_file(char* input_file);
void write_to_file(FILE* fp, char* to_write);
void sort_on_RAM(char* input_file, char* output_file);
void match_on_memory(char* i1, char* i2, char* o);
LIST* match_on_RAM(LIST* list1, LIST* list2);
void merge_on_memory(char* i1, char* i2, char* o);
LIST* merge_on_RAM(LIST* list1, LIST* list2);
```
