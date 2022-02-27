#!/bin/bash
if [ `which yum` ]; then
   yum install golang -y
   yum install python -y
elif [ `which apt` ]; then
   apt-get install golang -y
   apt-get install python -y
elif [ `which pacman` ]; then
   pacman -S go --noconfirm
   pacman -S python --noconfirm
elif [ `which dnf` ]; then
   dnf install golang -y
   dnf install python -y
else
   echo "Unknown OS"
fi