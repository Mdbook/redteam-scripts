# Killnuke

## Description

Simple shim to prevent the disabling of jumper. This project passes most commands to the actual killall binary, except if the `-o` or `-y` parameters are supplied, in which case it does nothing.

## Deployment

To install this project, run `killnuke.sh` as root.

## Dependencies

To run this project, you will need to have `golang` installed.
You can run `deploy/dependencies.sh` to install this automatically, or
run the commands below based on your distro:

Ubuntu/Debian
```bash
apt-get install golang-go
```
CentOS/RHEL
```bash
yum install golang
```
Arch
```bash
pacman -S go
```
