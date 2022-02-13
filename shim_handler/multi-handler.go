package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args
	for i := 1; i < len(args); i++ {
		ip := args[i]
		go do(ip)
	}
}

func do(ip string) {
	listenPort := strings.ReplaceAll(ip, ".", "")
	listenPort = "2" + listenPort[len(listenPort)-4:]
	cmd := exec.Command("xterm", "-title", ip, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
}
