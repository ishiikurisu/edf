#ifndef BUFFER_H
#define BUFFER_H 0
#include <stdlib.h>
#include <stdio.h>
#include "oa.h"

typedef struct
{
    char *source;
    char *mode;
    char *data;
    int size;
    bool available;
    FILE *stream;
}
BUFFER;

BUFFER* buffer_new(char *source, char *mode, int buffer_size)
{
    BUFFER *buffer = malloc(sizeof(BUFFER));

    buffer->source = source;
    buffer->mode = mode;
    buffer->data = string_new();
    buffer->size = buffer_size;
    buffer->stream = fopen(source, mode);
    buffer->available = 1;

    return buffer;
}

void buffer_status(BUFFER *buffer)
{
    printf("!%s\n", buffer->data);
}

BUFFER* buffer_flush(BUFFER *buffer)
{
    fwrite(buffer->data, strlen(buffer->data), 1, buffer->stream);
    free(buffer->data);
    buffer->data = string_new();
    return buffer;
}

BUFFER* buffer_feed(BUFFER *buffer)
{
    int char_count = 0;
    char read_char;

    while ((char_count < buffer->size) && (!feof(buffer->stream)))
    {
        read_char = fgetc(buffer->stream);
        cat(buffer->data, ctos(read_char));
        char_count++;
    }

    return buffer;
}

void buffer_write(BUFFER *buffer, char *to_write)
{
    cat(buffer->data, to_write);

    if (strlen(buffer->data) >= buffer->size)
        buffer = buffer_flush(buffer);
}

char buffer_read(BUFFER *buffer)
{
    char outlet = '\0';
    char *temp;

    if (strlen(buffer->data) == 0) {
        if (!feof(buffer->stream))
            buffer = buffer_feed(buffer);
        else {
            buffer->available = 0;
            return outlet;
        }
    }

    outlet = buffer->data[0];
    temp = substring(buffer->data, 1, strlen(buffer->data));
    free(buffer->data);
    buffer->data = temp;

    return outlet;
}

int buffer_is_available(BUFFER *buffer)
{
    return buffer->available;
}

int buffer_feof(BUFFER *buffer)
{
    return feof(buffer->stream);
}

char* buffer_read_line(BUFFER *buffer)
{
    char* line = string_new();
    char to_add = '\0';

    do {
        free(line);
        line = string_new();
        to_add = buffer_read(buffer);

        while (to_add != '\n' && to_add != '\0' && buffer->available)
        {
            cat(line, ctos(to_add));
            to_add = buffer_read(buffer);
        }
    }
    while (strlen(line) == 0 && buffer->available);

    if (strlen(line) == 0)
        free(line),
        line = NULL;

    return line;
}
char *buffer_readline(BUFFER *buffer)
{ return buffer_read_line(buffer); }

void buffer_close(BUFFER *buffer)
{
    if (buffer->mode[0] == 'a' || buffer->mode[0] == 'w')
        buffer_flush(buffer);
    buffer->available = 0;
    fclose(buffer->stream);
    free(buffer->data);
}

#endif
