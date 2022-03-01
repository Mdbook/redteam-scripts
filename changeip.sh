#Run this to change the host IP for every file in this repo.
read -p "Original IP: " ip1
echo
read -p "New IP: " ip2
find . \( ! -regex '.*/\..*' \) -type f | xargs sed -i 's#$ip1#$ip2#g'