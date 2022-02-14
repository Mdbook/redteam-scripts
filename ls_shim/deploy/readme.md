# Deploy

Tool to deploy ls_shim to all hosts on the same network.
This tool scans all devices on the 0/24 subnet and attempts to
install ls_shim on each of them by using credentials
specified by command line arguments.
## Dependencies

To run this project, you will need to have `golang` installed.
You can run `dependencies.sh` to install this automatically, or
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
## Usage/Examples

```
usage: go run deploy.go -u [username] -p [password] [args]
-v or --verbose                 |       Enable verbose output
-i [IPs] or --ignore [IPS]      |       Specify a list of IPs to ignore, separated by commas
-m or --multi                   |       Run in multithreaded mode. Not compatible with verbose.
-t [IP] or --target [IP]        |       Install on a remote machine & deploy from it
                                |       instead of the host machine. Not compatible with -i
--help or -h                    |       Display this help menu
--password-list [PASSWORDS]     |       Specify a list of passwords, separated by commas
--user-list [USERS]             |       Specify a list of users, separated by commas
```
## Using Target

Using the `-t` or `--target` option, you can specify this script to first install on a specified target,
and have that target deploy ls_shim to every device on its subnet. Please use with caution!