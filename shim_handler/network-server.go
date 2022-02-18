package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type connection struct {
	active chan bool
	ip     string
}

var connections []connection

func main() {
	for {
		GetPort()
	}
}

func setActive(ip string, active bool) {
	for i := 0; i < len(connections); i++ {
		if connections[i].ip == ip {
			connections[i].active <- active
			return
		}
	}
	conn := connection{make(chan bool), ip}
	connections = append(connections, conn)
}

func isActive(ip string) bool {
	for i := 0; i < len(connections); i++ {
		if connections[i].ip == ip {
			active := <-connections[i].active
			return active
		}
	}
	return false
}

func GetPort() {
	getPort, _ := net.Listen("tcp", "192.168.20.18:5003")
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)
	remoteIpForm := remoteIp[:strings.Index(remoteIp, ":")]
	remotePort := strings.ReplaceAll(remoteIpForm, ".", "")
	remotePort = "2" + remotePort[len(remotePort)-4:]
	if !isActive(remoteIpForm) {
		go do(remoteIpForm, remotePort)
		time.Sleep(100 * time.Millisecond)
	}
	conn.Write([]byte(remotePort))
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	return
}

func do(ip, listenPort string) {
	setActive(ip, true)
	cmd := exec.Command("xterm", "-title", ip, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
	setActive(ip, false)
}
