
# Service Herring

Installs up to 32 random services with random names, paths, and payloads.
To deploy this script, run the `install.sh` file. If you want to run the script
multiple times, you may run `./install.sh --demo` first and then run `./service-creator`
as many times as you'd like; however, this is not recommended. The `install.sh` file
automatically builds all payloads and the service-creator file, and then runs
service-creator, passing along any parameters.

Note: When deployed with the number of services set to the default amount (32), there will always be
at least one of each payload, with the exception of the `sleep` payload.

To view information about a specific payload, view the readme in its directory.

To deploy this project to an entire network of devices at the same time,
please view the readme in the `deploy` subfolder.


## Dependencies

To run this project, you will need to have `golang` installed.

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

Build using the bash script provided. This will build and automatically deploy service-creator; any parameters will be passed to the service-creator executable

```bash
#Default installation
sudo ./install.sh

#Build executables but don't deploy
sudo ./install.sh --demo

#Install and enable output for all payloads
sudo ./install.sh -v

```
    
## Usage/Examples
Any of these parameters may be used when running install.sh as well; they will be passed to service-creator when it is run by the script.

```bash
#Default run; create services and install them
./service-creator

#Create 20 services
./service-creator -n 20

#Generate services but do not install them
./service-creator --demo

#Enable verbose output for services
./service-creator -v

#Display help
./service-creator --help
./service-creator -h
```


## Documentation/Files modified

This script generates the following possible filenames:
```
b1ngus
borger
deleteme
dontdeleteme
file12345
flappy-bird-game
geck
heh
himom
homework
inconspicuous_file
issaservice
jeffUwU
lilboi
randomservice
temporary-file
top-secret
youfoundme
```
These files may be placed into any of the following  paths:
```
/etc/
/home/
/root/
/usr/lib/
/var/
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
