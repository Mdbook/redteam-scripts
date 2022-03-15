package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var port string
var hasPort bool = false
var takenPorts []string
var wg sync.WaitGroup

func main() {
	handleArgs()
	for {
		GetConnection()
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

func GetConnection() {
	getPort, _ := net.Listen("tcp", GetOutboundIP()+":"+port)
	conn, _ := getPort.Accept()
	defer conn.Close()
	defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s\n", remoteIp)
	// conn.Write([]byte("ls -al\n"))
	wg.Add(2)
	var active bool = true
	channel := make(chan string)
	go readConnection(conn, &active)
	go writeConnection(conn, &active, channel)
	wg.Wait()
	return
}

func readConnection(conn net.Conn, active *bool) {
	defer wg.Done()
	for {
		buf := make([]byte, 65535)
		_, err := conn.Read(buf)
		if err != nil {
			*active = false
			fmt.Println(err)
			break
		}
		fmt.Println(string(buf))
	}
}

func writeConnection(conn net.Conn, active *bool, channel chan string) {
	defer wg.Done()
	for {
		if *active {
			text := <-channel
			conn.Write([]byte(fmt.Sprintf("%s\n", text)))
		} else {
			return
		}

	}

}

func readStdin(channel chan string) {
	channel <- "hello"
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	text, _ := reader.ReadString('\n')
	// }
}
