#!/bin/bash

#Build payloads
cd downloader
go build downloader.go
cd ../file-creator
go build file-creator.go
cd ../random-messenger
go build random-messenger.go
cd ../shell
go build shell.go
cd ../user-creator
go build user-creator.go
cd ..
go build service-creator.go
echo "Built go files"