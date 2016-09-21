#pragma once
#include <oa.h>
#include <buffer.h>

char *set_name(char *input)
{
	LIST *bits = list_strsplit(input, '.');
	LIST *bit = bits->next;
	char *out = bit->value;

	for (inc(bit); bit->next != NULL; inc(bit))
	{
		cat(out, ctos('.'));
		ape(out, bit->value);
	}

	cat(out, ".ascii");
	list_free(bits);
	return out;
}

char *treat_line(char *line)
{
	LIST *values = list_strsplit(line, ',');
	char *output = string_new();
	LIST *it = NULL;

	for (it = values->next; it != NULL; inc(it))
	{
		ape(output, it->value);
		cat(output, ctos(' '));
	}

	cat(output, ctos('\n'));
	list_free(values);
	return output;
}

void csv2fullascii(char *input)
{
	char *output = NULL;
	BUFFER *csv = NULL;
	BUFFER *ascii = NULL;
	char *line = NULL;

	output = set_name(input);
	csv = buffer_new(input, "r", 2048);
	ascii = buffer_new(output, "w", 2048);

	buffer_readline(csv);
	line = buffer_readline(csv);
	while (buffer_is_available(csv))
	{
		buffer_write(ascii, treat_line(line));
		free(line);
		line = buffer_readline(csv);
	}

	buffer_close(csv);
	buffer_close(ascii);
}
