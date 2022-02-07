package main

import (
	"log"
	"net"
	"strings"

	//"strconv"
	"fmt"
	//"bufio"
)

func main() {
	for {
		GetPort()
	}

	//fmt.Printf(string(GetOutboundIP()))

}

func GetPort() {
	getPort, _ := net.Listen("tcp", "192.168.12.6:5003")
	//defer getPort.Close()
	conn, _ := getPort.Accept()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)

	remoteIpForm := remoteIp[:strings.Index(remoteIp, ":")]
	remotePort := strings.ReplaceAll(remoteIpForm, ".", "")
	remotePort = "2" + remotePort[len(remotePort)-4:]
	conn.Write([]byte(remotePort))
	//remotePortInt, _ := strconv.Atoi(remotePort)
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	conn.Close()
	getPort.Close()
	return
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ip := localAddr.IP
	ipstr := ip.String()
	ipstr = strings.ReplaceAll(ipstr, ".", "")
	return ipstr
}
