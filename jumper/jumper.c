
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <sys/time.h>
#include <time.h>

#define FALSE 0
#define TRUE 1
#define ERROR -1
#define SLEEPTIME 1;

int isKill(){
    if (access("stop", F_OK) == FALSE ) {
        //file exists
        return TRUE;
    }
    //file does not exist
    return FALSE;
}


void rand_str(char *dest, size_t length) {
    struct timeval te;
    gettimeofday(&te, NULL);
    srand((unsigned int)te.tv_usec);
    char charset[] = "0123456789"
                     "abcdefghijklmnopqrstuvwxyz"
                     "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

    while (length-- > 0) {
        size_t index = (double) rand() / RAND_MAX * (sizeof charset - 1);
        *dest++ = charset[index];
    }
    *dest = '\0';
}


int main(int argc, char *argv[]){
    setuid(0);
    if (isKill() == TRUE){
        return 1;
    }
    int string_len = 5;
    char str[] = { [5] = '\1' }; // make the last character non-zero so we can test based on it later
    rand_str(str, sizeof str - 1);
    assert(str[5] == '\0');      // test the correct insertion of string terminator
    printf("%s\n", str);
    
}

