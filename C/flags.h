#ifndef FLAGS_H
#define FLAGS_H
#include <stdlib.h>
#include <signal.h>
#include <oa.h>

typedef struct {
    char *input_file;
    bool single;
} 
OPTIONS;

enum FLAG 
{
    SINGLE,
    MULTIPLE
};

enum FLAG which_flag(char *arg)
{
    enum FLAG flag = SINGLE;

    if (compare(arg, "-m") == EQUAL || compare(arg, "--multiple") == EQUAL)
        flag = MULTIPLE;
    else if (compare(arg, "-s") == EQUAL || compare(arg, "--single") == EQUAL)
        flag = SINGLE;
    else
        printf("Invalid flag\n"),
        raise(SIGILL);

    return flag;
}

bool is_flag(char *arg)
{
    return (arg[0] == '-')? true : false;
}

OPTIONS* parse_flags(int argc, char **argv)
{
    OPTIONS *options = (OPTIONS*) malloc(sizeof(OPTIONS));
    int i = 0;

    options->single = true;
    options->input_file = NULL;

    for (i = 1; i < argc; ++i)
    {
        if (is_flag(argv[i])){
            if (which_flag(argv[i]) == MULTIPLE)
                options->single = false;
        }
        else {
            if (options->input_file == NULL)
                options->input_file = argv[i];
            else
                printf("Select only one file\n"),
                raise(SIGILL);
        }
    }

    return options;
}

#endif
