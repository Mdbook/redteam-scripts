package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var port string
var hasPort bool = false
var takenPorts []string

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

func getRandomPort() string {
	remotePort := "2" + strconv.Itoa(random(99)) + strconv.Itoa(random(99))
	if findIndex(takenPorts, remotePort) == -1 {
		return remotePort
	}
	return getRandomPort()
}

func GetPort() {
	getPort, _ := net.Listen("tcp", "10.1.1.6:"+port)
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)
	remoteIpForm := remoteIp[:strings.Index(remoteIp, ":")]
	remotePort := getRandomPort()
	takenPorts = append(takenPorts, remotePort)
	// remotePort := strings.ReplaceAll(remoteIpForm, ".", "")
	// remotePort = "2" + remotePort[len(remotePort)-4:]
	go do(remoteIpForm, remotePort)
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte(remotePort))
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	return
}

func do(ip, listenPort string) {
	fmt.Println(takenPorts)
	cmd := exec.Command("xterm", "-title", ip+" | "+port, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
	takenPorts, _ = remove(takenPorts, findIndex(takenPorts, listenPort))
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func remove(slice []string, i int) ([]string, string) {
	//Remove an item from a slice
	name := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, name

}

func findIndex(slice []string, value string) int {
	//Find the index of a string in a slice
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}
