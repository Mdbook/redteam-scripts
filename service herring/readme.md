
# Service Herring

A brief description of what this project does and who it's for


## Dependencies

To run this project, you will need to have `golang` installed.

Ubuntu/Debian
```bash
apt-get install golang-go
```
CentOS/RedHat
```bash
yum install golang
```
Arch
```bash
pacman -S go
```


## Installation

Build using the bash script provided

```bash
sudo ./setup.sh
```
    
## Usage/Examples

```bash
#Default run; create services and install them
./service-creator

#Create 20 services
./service-creator -n 20

#Generate services but do not install them
./service-creator --demo

#Display help
./service-creator --help
./service-creator -h
```

