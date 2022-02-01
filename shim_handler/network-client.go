package main

import (
    "net"
	"strings"
	"log"
	"fmt"
	"bufio"
	"strconv"
    )


func main () {
	num := GetPort()
	fmt.Printf("Port is: %d\n", num)
	//fmt.Printf(string(GetOutboundIP()))

}

func GetPort() int {
	getPort,err := net.Dial("tcp", "192.168.3.6:5003")
	if err != nil {
		fmt.Println("Couldn't get connection")
		return 420
	}
	defer getPort.Close()
	status, _ := bufio.NewReader(getPort).ReadString('\n')
	num, _ := strconv.Atoi(status)
	return num
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