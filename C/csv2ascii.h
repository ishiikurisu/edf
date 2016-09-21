#ifndef CSV_TO_ASCII
#define CSV_TO_ASCII
#include "stdio.h"
#include "stdlib.h"
#include "stdbool.h"
#ifndef BUFFER_H
#include "buffer.h"
#endif

/*
# Common functions

+ void write_line(BUFFER *outlet, LIST *line)
+ int get_chan(char* line)
+ LIST* equal_split(char* labels, int chan)
+ char* get_labels(char* line)
+ LIST* parse_header(char *line);
+ LIST* parse_line(char *line);
*/

char* are_these_labels(char *field)
{
    char *clean = tidy_string(field);
    char *temp = NULL;

    if (match("labels", clean))
        temp = substring(clean, 7, strlen(clean)),
        free(clean),
        clean = temp;
    else
        free(clean),
        clean = NULL;

    return clean;
}

bool is_label_valid(char *label)
{
    static bool ext = false;

    if (compare(label, "EEG"))
        return false;
    if (compare(label, "EXT")) {
        if (ext == false) {
            ext = true;
            return false;
        }
        else {
            return true;
        }
    }

    return true;
}

LIST* filter_annotations(LIST* stuff)
{
    int index = 0;
    LIST* it = NULL;

    for (it = stuff->next; it != NULL; inc(it))
        if (match(it->value, "EDF Annotations"))
            break;
        else
            index++;
    stuff = list_remove(stuff, index);
    for (it = stuff->next; it != NULL; inc(it))
        it->value = tidy_string(it->value);

    return stuff;
}

/**
* Gets the number of channels in this CSV file
* @param line the header
*/
int get_chan(char *line)
{
    LIST *fields = NULL;
    LIST *field = NULL;
    char *key = NULL;
    int result = -1;

    if (strlen(line) <= 1) return result;

    fields = list_strsplit(line, ',');
    for (field = fields->next; field != NULL && key == NULL; inc(field))
        if (match("chan", tidy_string(field->value)))
            key = substring(tidy_string(field->value), 5, strlen(tidy_string(field->value)));
        else
            key = NULL;

    list_free(fields);
    sscanf(key, "%d", &result);
    return result;
}

/**
* Splits a string into equal pieces
* @param labels the part of the string containing labels
* @param chan the number of channels
* @return the list of labels
*/
LIST* equal_split(char *labels, int chan)
{
    LIST *outlet = list_new();
    char *temp = NULL;
    int total = strlen(labels);
    int piece = total / chan;
    int i, j;

    for (j = 0; j < chan; ++j)
        outlet = list_add(outlet, substring(labels, j * piece + j, (j+1) * piece + j));

    return outlet;
}

/**
 * Extract the 'labels' field as a list from the CSV header
 * @param  line   a c_string containing the file's header
 * @return fields a joe_list containing the label for each signal
 */
char* get_labels(char *line)
{
    LIST *fields = NULL;
    LIST *field = NULL;
    char *labels = NULL;

    if (strlen(line) <= 1) return NULL;
    fields = list_strsplit(line, ',');
    for (field = fields->next; field != NULL && labels == NULL; inc(field))
        labels = are_these_labels(field->value);

    list_free(fields);
    return labels;
}

/**
 * Extract the 'labels' field as a list from the CSV header as a list
 * @param  line   a c_string containing the file's header
 * @return fields a joe_list containing the label for each signal
 */
LIST* parse_header(char *line)
{
    char *labels = get_labels(line);
    LIST* stuff = equal_split(labels, get_chan(line));
    // stuff = list_strsplit(labels, ' ');

    stuff = filter_annotations(stuff);

    return stuff;
}

/**
 * Separates the values in a CSV line
 * @param  line   a CSV file line in c_str format
 * @return values a joe_string containing each value in
 */
LIST *parse_line(char *line)
{
    LIST *values = NULL;
    LIST *value = NULL;
    char *temp = NULL;

    if (line == NULL || strlen(line) <= 1)
        return NULL;

    values = list_strsplit(line, ',');
    for (value = values->next; value != NULL; inc(value))
        temp = tidy_string(value->value),
        free(value->value),
        value->value = temp;

    return values;
}

/**
 * writes a new ASCII line in the selected buffer
 * @param outlet joe_buffer for the output file
 * @param stuff  joe_list containing the CSV columns
 */
void write_line(BUFFER *outlet, LIST *stuff)
{
    LIST *item = stuff->next;

    buffer_write(outlet, item->value);
    for (inc(item); item != NULL; inc(item))
        buffer_write(outlet, concat(ctos(' '), item->value));
    buffer_write(outlet, "\n");
}

/*
# Main functions

+ void csv2single(char *input);
+ void csv2multiple(char *input);
*/

/**
 * Translates a `csv` file to the ascii format
 * @param input  csv file name
 * @param output ascii file name
 */
#include "csv2ascii/single.h"
void csv2single(char *input)
{
    BUFFER *inlet = buffer_new(input, "r", 256);
    BUFFER *outlet = buffer_new(single_get_output(input), "w", 256);
    LIST *line = parse_header(buffer_readline(inlet));

    while ((line = parse_line(buffer_readline(inlet))) != NULL)
        write_line(outlet, line),
        list_free(line);

    buffer_close(inlet);
    buffer_close(outlet);
}

/**
 * separates a `csv` file into many `ascii` files, one for each signal.
 * @param input a c_string naming the `csv` file
 */
#include "csv2ascii/multiple.h"
void csv2multiple(char *input)
{
    BUFFER *inlet = buffer_new(input, "r", 256);
    LIST *line = parse_header(buffer_readline(inlet));
    BUFFER **outlets = multiple_buffers_new(input, line);

    /* TODO: add case where there are more channel names than expected */
    for (line = parse_line(buffer_readline(inlet));
         line != NULL;
         line = parse_line(buffer_readline(inlet)))
        multiple_write_lines(outlets, line),
        list_free(line);

    buffer_close(inlet);
    multiple_buffers_close(outlets);
}

#endif
