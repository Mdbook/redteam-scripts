package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	for {
		GetPort()
	}
}

func GetPort() {
	getPort, _ := net.Listen("tcp", "192.168.20.18:5003")
	//defer getPort.Close()
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)

	remoteIpForm := remoteIp[:strings.Index(remoteIp, ":")]
	remotePort := strings.ReplaceAll(remoteIpForm, ".", "")
	remotePort = "2" + remotePort[len(remotePort)-4:]
	fmt.Println(remotePort)
	//go do(remoteIpForm, remotePort)
	fmt.Println("here")
	conn.Write([]byte(remotePort))
	//remotePortInt, _ := strconv.Atoi(remotePort)
	fmt.Printf("Sent port %s to %s\n\n", remotePort, remoteIp)
	return
}

func do(ip, listenPort string) {
	//defer wg.Done()
	cmd := exec.Command("xterm", "-title", ip, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()

}

/*
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
*/
