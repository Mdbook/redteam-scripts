
# Downloader

Fetches commands from a HTTP server and executes them on the local machine.


## Usage/Examples

This program reads instructions from a file named `stat` located in the
root directory of a webserver. The first line of this file must contain only the word `EDITION`,
and then a number from 0-99. This number should be increased every time the host wants to deploy a new
set of commands. After the first line, any of the following commands may be put in the file:


#### DOWNLOAD
Downloads a file from the host server, with the path to the filename located on the line
directly below.

For example, the following line will download the file located at `[HOST]/files/foo.txt`:
```
DOWNLOAD
files/foo.txt
```

#### EXECUTE
Executes a binary on the system. The path to the file must be included on the next line.
The filename cannot contain spaces.

Example:
```
EXECUTE
/path/to/file
```

#### RUN
Executes the shell command found on the next line as root.

Exammples:
```
RUN
echo "Somebody is in your system..." > /home/foo/readme.txt

RUN
iptables -F && cp /home/foo/readme.txt /tmp/
```

#### MOVE
Moves a file from one location to another. The source path must be put on the line after
the command, and the destination file on the line after that.

Example:
```
MOVE
/home/foo/readme.txt
/tmp/dontreadme.txt
```
This will move `/home/foo/readme.txt` to `/tmp/dontreadme.txt`


#### COPY
Copies a file to another location. The source path must be put on the line after
the command, and the destination file on the line after that.

Example:
```
COPY
/home/foo/readme.txt
/tmp/dontreadme.txt
```
This will copy `/home/foo/readme.txt` to `/tmp/dontreadme.txt`


#### DELETE
Deletes the file specified on the next line.

Example:
```
DELETE
/tmp/dontreadme.txt
```

#### MESSAGE
Sends a message to all users using the wall command.
The message to be sent must be put on the next line.

Example:
```
MESSAGE
hello there
```

#### SLEEP
Pauses execution of the script for the amount of time
specified in the next line (in seconds).

Example:
```
SLEEP
10
```

## Parameters

`-v`: enable output

Example:
```
go run downloader.go -v

go build downloader.go && ./downloader -v
```