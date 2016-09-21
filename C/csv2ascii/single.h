#ifdef CSV_TO_ASCII
#ifndef SINGLE_H
#define SINGLE_H 0

char* single_get_output(char *input)
{
    char *output = substring(input, 0, strlen(input)-4);
    cat(output, ".ascii");
    return output;
}

#endif
#endif