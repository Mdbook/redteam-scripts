#!/bin/bash
cp master-server.go vi-server.go
cp master-server.go vim-server.go
cp master-server.go nano-server.go
cp master-server.go ls-server.go

sed -i 's/{SERVERNAME}/LS/g' ls-server.go
sed -i 's/{SERVERNAME}/VI/g' vi-server.go
sed -i 's/{SERVERNAME}/VIM/g' vim-server.go
sed -i 's/{SERVERNAME}/NANO/g' nano-server.go

sed -i 's/{SERVERPORT}/5003/g' ls-server.go
sed -i 's/{SERVERPORT}/5004/g' vi-server.go
sed -i 's/{SERVERPORT}/5005/g' vim-server.go
sed -i 's/{SERVERPORT}/5006/g' nano-server.go

sed -i 's/{ASSIGNEDPORT}/2/g' ls-server.go
sed -i 's/{ASSIGNEDPORT}/3/g' vi-server.go
sed -i 's/{ASSIGNEDPORT}/4/g' vim-server.go
sed -i 's/{ASSIGNEDPORT}/5/g' nano-server.go

if [ $1 ]; then
    if [ $1 == "--nostart" ]; then
        exit
    fi
fi

go run run-helper.go

echo "Processes started."

sleep 0.2
rm ls-server.go
rm vi-server.go
rm vim-server.go
rm nano-server.go