package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var port string
var hasPort bool = false

func main() {
	handleArgs()
	for {
		GetPort()
	}
}

func handleArgs() {
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-p":
				port = args[i+1]
				hasPort = true
			}
		}
	}
	if !hasPort {
		fmt.Printf("Port: ")
		fmt.Scanln(&port)
		hasPort = true
	}
	fmt.Println("Listening on port " + port)
}

func GetPort() {
	getPort, _ := net.Listen("tcp", "192.168.20.18:"+port)
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)
	remoteIpForm := remoteIp[:strings.Index(remoteIp, ":")]
	remotePort := strings.ReplaceAll(remoteIpForm, ".", "")
	remotePort = "2" + remotePort[len(remotePort)-4:]
	go do(remoteIpForm, remotePort)
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte(remotePort))
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	return
}

func do(ip, listenPort string) {
	cmd := exec.Command("xterm", "-title", ip+" | "+port, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
}
