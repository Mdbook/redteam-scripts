# Shell
Listens on a randomly selected tcp port between `62000-62999`.
Any user who connects to the port is given a root shell.
## Parameters

`-v`: enable verbose

Examples:
```
go run shell.go -v

go build shell.go && ./shell
```