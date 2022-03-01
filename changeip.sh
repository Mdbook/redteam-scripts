#Run this to change the host IP for every file in this repo.
read -p "Original IP: " ip1
read -p "New IP: " ip2
find . -type f | xargs sed -i  "s/$ip1/$ip2/g"
echo "IP changed."