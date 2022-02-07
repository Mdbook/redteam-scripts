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
cd ../sleep
go build sleep.go
echo "Built sleep"
cd ..
go build service-creator.go
echo "Built service-creator"

echo "Built all go files"
echo
echo

echo "Running service-creator..."
echo
./service-creator $@