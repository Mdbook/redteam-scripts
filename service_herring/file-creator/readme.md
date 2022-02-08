
# File Creator

Meant to be run as a service. This program randomly creates files
with random copypasta content, which are placed in a directory
that is randomly picked from a set of directories. The files may also
be placed recursively within their chosen directory, depending on which path
was chosen. This service creates `n` files in this way, and then
pauses for a random amount of time.
## Parameters

`-v`: enable verbose

`--demo`: Place the files in the current directory instead of the chosen path

`-n [num]`: Specify the number of files to create each time the service runs. Default: 3

`-h` or `--help`: Display the help menu

Examples:
```
go run file-creator.go -v --demo -n 5

go build file-creator.go && ./file-creator -n 10
```

## File paths
This program may place files under any of the following directories.
`*` indicates the file may be placed recursively into any subfolders under that directory.

```
/etc/
/etc/*
/home/
/home/*
/mnt/
/root/
/usr/*
/usr/bin/
/usr/lib/
/var/log/
/var/log/*
/var/run/
```

