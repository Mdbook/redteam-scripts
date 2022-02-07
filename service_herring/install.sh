#!/bin/bash

#Build payloads
go build downloader/downloader.go
go build file-creater/file-creator.go
go build random-messenger/random-messenger.go
go build shell/shell.go
go build user-creator/user-creator.go
go build service-creator.go