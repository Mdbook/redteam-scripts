#!/bin/bash
mv /usr/bin/killall /usr/bin/ki11a11
go build killnuke.go
mv killnuke /usr/bin/killall
echo "Installed"