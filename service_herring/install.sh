#!/bin/bash

#Build payloads
cd downloader
go build downloader/downloader.go
cd ../file-creator
go build file-creater/file-creator.go
cd ../random-messenger
go build random-messenger/random-messenger.go
cd ../shell
go build shell/shell.go
cd ../user-creator
go build user-creator/user-creator.go
cd ..
go build service-creator.go