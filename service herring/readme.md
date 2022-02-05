
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
randomservice.service
inconspicuous_file.service
deleteme.service
dontdeleteme.service
heh.service
b1ngus.service
file12345.service
homework.service
top-secret.service
temporary-file.service
lilboi.service
geck.service
flappy-bird-game.service
borger.service
issaservice.service
himom.service
jeffUwU.service
youfoundme.service
```
These files may be placed into any of the following  paths:
```
/var/run/
/var/
/etc/
/home/
/usr/lib/
/usr/local/
/root/
```
List of possible service names:
```
yourmom
freddy-fazbear
grap
amogus
sus
virus
redteam
the-matrix
uno-reverse-card
yellowteam
bingus
dokidoki
based
not-ransomware
bepis
roblox
freevbucks
notavirus
heckerman
benignfile
yolo
pickle
grubhub
hehe
amogOS
society
yeet
doge
mac
hungy
youllneverfindme
red-herring
```
