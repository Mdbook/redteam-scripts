#!/bin/bash
if [ `which yum` ]; then
   yum install golang -y
   yum install libcurl-devel -y
elif [ `which apt-get` ]; then
   apt-get install golang -y
   apt-get install libcurl4-openssl-dev -y
elif [ `which pacman` ]; then
   pacman -S go --noconfirm
   pacman -S curl --noconfirm
elif [ `which dnf` ]; then
   dnf install golang -y
   dnf install libcurl-devel -y
else
   echo "Unknown OS"
fi