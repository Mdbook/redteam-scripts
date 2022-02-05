
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


## Documentation/Files modified

This script generates the following possible filenames:
```
b1ngus.service
borger.service
deleteme.service
dontdeleteme.service
file12345.service
flappy-bird-game.service
geck.service
heh.service
himom.service
homework.service
inconspicuous_file.service
issaservice.service
jeffUwU.service
lilboi.service
randomservice.service
temporary-file.service
top-secret.service
youfoundme.service
```
These files may be placed into any of the following  paths:
```
/etc/
/home/
/root/
/usr/lib/
/usr/local/
/var/
/var/run/
```
List of possible service names:
```
amogOS
amogus
based
benignfile
bepis
bingus
doge
dokidoki
freddy-fazbear
freevbucks
grap
grubhub
heckerman
hehe
hungy
mac
not-ransomware
notavirus
pickle
red-herring
redteam
roblox
society
sus
the-matrix
uno-reverse-card
virus
yeet
yellowteam
yolo
youllneverfindme
yourmom
```
