#!/bin/bash
if [ `which yum` ]; then
   yum install golang -y
elif [ `which apt` ]; then
   apt-get install golang -y
elif [ `which pacman` ]; then
   pacman -S go --noconfirm
elif [ `which dnf` ]; then
   dnf install golang -y
   dnf install openssl -y
else
   echo "Unknown OS"
fi