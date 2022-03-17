package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	for {
		do()
	}

}
func do() {
	listener, _ := net.Listen("tcp", GetOutboundIP()+":8003")
	conn, _ := listener.Accept()
	defer listener.Close()
	defer conn.Close()
	conn.Write([]byte("dir\n"))
	fmt.Println("sent")
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Print(str)
	}
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
