#ifdef CSV_TO_ASCII
#ifndef MULTIPLE_H
#define MULTIPLE_H 1


char* _multiple_get_surname(char* input)
{
	char *output = string_new();
	LIST *bits = list_strsplit(input, '.');
	LIST *bit = bits->next;

	cat(output, bit->value);
	for (inc(bit); bit->next != NULL; inc(bit))
		cat(output, ctos('.')),
		cat(output, bit->value);

	free(bits);
	return output;
}

/**
 * Create many buffers to files named as their mother files 
 * and their respective files
 * @param  input   the name of the mother 
 * @param  line    a joe_list containing the labels on the file
 * @return buffers a NULL-ended array with buffers to the files
 */
BUFFER** multiple_buffers_new(char *input, LIST *line)
{
	BUFFER **outlets = (BUFFER**) malloc(sizeof(BUFFER*) * (list_length(line)+1));
	char *surname = _multiple_get_surname(input);
	char *name = NULL;
	LIST *it = NULL;
	int i = 0;

	for (it = line->next, i = 0; it != NULL; inc(it), ++i)
		name = concat(concat(surname, concat(ctos('.'), it->value)), ".ascii"),
		outlets[i] = buffer_new(name, "w", 64),
		free(name);

	outlets[i] = NULL;
	free(surname);
	return outlets;
}

/**
 * writes a joe_list to the buffers following their order
 * @param outlets a NULL-ended array of joe_buffer pointers
 * @param line    a joe_list containg the data to be written
 */
void multiple_write_lines(BUFFER **outlets, LIST *line)
{
	BUFFER** buffer = NULL;
	LIST* it = NULL;

	for (buffer = outlets, it = line->next; (*buffer) != NULL; ++buffer, inc(it))
		buffer_write((*buffer), concat(it->value, ctos('\n')));
}

/**
 * closes many buffers at once
 * @param buffers a NULL-ended array of joe_buffer pointers
 */
void multiple_buffers_close(BUFFER **buffers)
{
	while ((*buffers))
		buffer_close((*buffers)),
		++buffers;
}

#endif
#endif