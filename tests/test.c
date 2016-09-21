#include "stdio.h"
#include "stdlib.h"
#include "C/buffer.h"
#include "csv2ascii.h"

int main(int argc, char *argv[]) {
    BUFFER *inlet = buffer_new(argv[1], "r", 256);
    char* header = buffer_readline(inlet);
    int chan = 0;
    char *labels = NULL;
    LIST *stuff = NULL;

    printf("# Getting channels\n");
    printf("---\n");
    chan = get_chan(header);
    labels = get_labels(header);
    printf("chan: %d\n", chan);
    printf("labels: %s\n", labels);
    printf("\n");
    printf("# Equally splitting\n");
    // char* list_yaml(LIST*)
    stuff = equal_split(labels, chan);
    printf("%s\n", list_yaml(stuff));
    printf("# Filtering annotations\n");
    printf("%s\n", list_yaml(filter_annotations(stuff)));

    printf("...\n");
    buffer_close(inlet);
    return 0;
}
