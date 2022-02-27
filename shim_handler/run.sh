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

if [ $1 ]; then
    if [ $1 == "--nostart" ]; then
        exit
    fi
fi
xterm -title "LS | MASTER" -e "go run ls-server.go" | &
xterm -title "VI | MASTER" -e "go run vi-server.go" | &
xterm -title "VIM | MASTER" -e "go run vim-server.go" | &
xterm -title "NANO | MASTER" -e "go run nano-server.go" | &
# TODO: fix this
echo "Processes started."

sleep 0.2
rm ls-server.go
rm vi-server.go
rm vim-server.go
rm nano-server.go