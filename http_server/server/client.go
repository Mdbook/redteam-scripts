package main

import (
	"fmt"
	"net"
)

type Client struct {
	id         int
	wanIP      string
	lanIP      string
	port       string
	clientType string
	os         string
	osFlavor   string
	isEncoded  bool
	conn       net.Conn
}

func handleClient(clientInfo ClientInfo, port string) {
	getPort, _ := net.Listen("tcp", GetOutboundIP()+":"+port)
	conn, _ := getPort.Accept()
	// defer conn.Close()
	// defer getPort.Close()
	client := globalMap.CreateClient(clientInfo, port, conn)
	fmt.Printf("New client (%s) connected with id: %d \n\n", clientInfo.lanIP, client.id)
	go runReadClient(client)
	go runWriteClient(client)
}

func runReadClient(client Client) {
	conn := client.conn
	// conn.Write([]byte(fmt.Sprintf("%s\n", "dir")))
	for {
		if client.clientType == "Encoded Reverse Shell" {
			// TODO: handle encoded reverse shell
		}
		// TODO: base64 encode
		// buf, err := bufio.NewReader(conn).ReadString('\n')
		buf := make([]byte, 65535)
		_, err := conn.Read(buf)
		if err != nil {
			// TODO error handling
			fmt.Println(err)
			break
		}
		if globalMap.GetActiveChannel() == client.id {
			fmt.Print(string(buf))
			// if globalMap.IsSingle() {
			// fmt.Println()
			// }
		}
	}
}

func runWriteClient(client Client) {
	conn := client.conn
	channel := globalMap.GetChannel(client.id)
	isActive := <-*channel
	for {
		if isActive && globalMap.GetActiveChannel() == client.id {
			stdin := globalMap.GetStdin()
			stdReadLine := <-*stdin
			if stdReadLine == "!!!FIN!!!" {
				isActive = false
			} else {
				conn.Write([]byte(fmt.Sprintf("%s\n", trim(stdReadLine))))
			}
		} else {
			isActive = <-*channel
		}
	}

}
