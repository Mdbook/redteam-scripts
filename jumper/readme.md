# Jumper

Program that relocates itself every ~50-150 milliseconds, executing a payload each time it jumps.

The jumper may move to any of the following locations:

```
/usr/local/
/usr/bin/
/etc/
/var/lib/
/usr/sbin/
/root/
/home/
```

## Usage

Any payload to be executed should be put into the `payload()` function in `payload.c`. After programming the payload, simply compile `jumper.c` and execute it on the host machine.

For windows, use `windows-jumper.c` and `windows-payload.c` respectively.

## Default Payload

The default payload has two functions: `stopServices()` and `dropAllFirewall()`.

`stopServices` stops common scored blue team services: `sshd`, `nginx`, `apache2`, as well as disabling icmp.

`dropAllFirewall` clears all iptables rules, then sets the default policy to `DROP` for all tables.


## Default Payload- Windows

The default payload for windows downloads a background from `server.mdbooktech.com` if the file doesn't already exist, and then sets the current background to the image if the current background is not already that image.

## SaveCPU mode

The linux version of jumper has a define statement called `SAVECPU`. If set to `1` before compiling, SaveCPU mode will be enabled. This will make the payload execute every 9th jump. If `SAVECPU` is set to `0`, the payload will instead be executed every jump.