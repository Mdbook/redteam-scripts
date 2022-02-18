#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main (int argc, char *argv[]) {
    char args[100] = "vi ";
    for (int i = 1;i<argc;i++){
        strcat(args, argv[i]);
        if (i != argc-1){
            strcat(args, " ");
        }
    }
    // system("xterm &");
    // system(args);
    pid_t fork_pid = fork();
    if (fork_pid == 0) {

        printf("Hello from the child!\n");
    } else {
        printf("Hello from the parent!\n");
    }
}
