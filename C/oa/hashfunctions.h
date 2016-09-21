#include <math.h>
#include <limits.h>

int stupid(char* input)
{
    int output = 0;
    int limit = strlen(input);
    int i;

    for (i = 0; i < limit; ++i)
        output += input[i] * (i + 1);

    return output;
}

int dumb(char* input)
{
    int output = 0;
    int limit = strlen(input);
    int i;

    for (i = 0; i < limit; ++i)
        output += input[i];

    return (output << 2) + (output >> 2);
}

int knuth(char* input)
{
    int result = 0;
    int limit = strlen(input);
    int i = 0;

    for (i = 0; i < limit; ++i)
        result = (result + input[i] * 7) % (INT_MAX);
    result *= 3557;
    result = abs(result);

    return result;
}

int MD5(char* input)
{
    int result = 0;
    int limit = strlen(input);
    int constant = 0;
    int i = 0;

    for (i = 0; i < limit; ++i)
    {
        constant += (int) abs(sin(i + 1) * INT_MAX);
        result = result * 10 + input[i] * constant;
    }

    return result;
}
