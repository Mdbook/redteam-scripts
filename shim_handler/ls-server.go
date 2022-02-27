package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

var port string = "5003"
var hostIP string = "192.168.20.18"
var hasPort bool = false

func main() {
	handleArgs()
	for {
		GetPort()
	}
}

func handleArgs() {
	fmt.Println("Listening on port " + port)
}

func GetPort() {
	getPort, _ := net.Listen("tcp", hostIP+":"+port)
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
	cmd := exec.Command("xterm", "-title", ip+" (LS)", "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
}
