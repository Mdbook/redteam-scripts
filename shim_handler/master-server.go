package main

// TODO: look at this code
// TODO: Make sure the IAP you get is the LAN ip and not the WAN ip
import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var port string = "{SERVERPORT}"
var portPrefix string = "{ASSIGNEDPORT}"
var hostIP string = "10.100.1.101"
var hasPort bool = false
var takenPorts []string

func main() {
	hostIP = GetOutboundIP()
	handleArgs()
	for {
		GetPort()
	}
}

func handleArgs() {
	fmt.Println("Listening on port " + port)
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

func getRandomPort() string {
	port1 := strconv.Itoa(random(10))
	port2 := strconv.Itoa(random(99))
	if len(port2) <= 1 {
		port2 = "0" + port2
	}
	// fmt.Println(port1)
	// fmt.Println(port2)
	remotePort := "2" + portPrefix + port1 + port2
	if findIndex(takenPorts, remotePort) == -1 {
		return remotePort
	}
	return getRandomPort()
}

func GetPort() {
	getPort, _ := net.Listen("tcp", hostIP+":"+port)
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp, err := bufio.NewReader(conn).ReadString('\n')
	remoteIpForm := remoteIp
	if err != nil {
		fmt.Println(err.Error())
	}
	if remoteIp == "none\n" || remoteIp == "none" {
		remoteIp = conn.RemoteAddr().String()
		remoteIpForm = remoteIp[:strings.Index(remoteIp, ":")]
		remoteIp = remoteIpForm
	}
	fmt.Printf("Received request from %s\n", remoteIp)
	remotePort := getRandomPort()
	takenPorts = append(takenPorts, remotePort)
	go do(remoteIpForm, remotePort)
	time.Sleep(200 * time.Millisecond)
	conn.Write([]byte(remotePort))
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	time.Sleep(5 * time.Second)
	return
}

func do(ip, listenPort string) {
	cmd := exec.Command("xterm", "-title", ip+" ({SERVERNAME})", "-e", "nc", "-l", "-p", listenPort)
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
