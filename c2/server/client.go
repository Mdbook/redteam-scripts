package main

import (
	"bufio"
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
	hostname   string
	isEncoded  bool
	isDead     bool
	conn       net.Conn
	listener   net.Listener
}

type ClientInfo struct {
	lanIP      string
	clientType string
	os         string
	osFlavor   string
	hostname   string
	isEncoded  bool
}

func clientDisconnect(client Client) {
	client.conn.Close()
	client.listener.Close()
	globalMap.KillClient(client.id)
}

func handleClient(clientInfo ClientInfo, port string) {
	getPort, _ := net.Listen("tcp", GetOutboundIP()+":"+port)
	conn, _ := getPort.Accept()
	// defer conn.Close()
	// defer getPort.Close()
	client := globalMap.CreateClient(clientInfo, port, conn, getPort)
	fmt.Printf("New client (%s) connected with id: %d \n", clientInfo.lanIP, client.id)
	caret()
	go runReadClient(client)
	go runWriteClient(client)
}

func runReadClient(client Client) {
	conn := client.conn
	for {
		if client.isEncoded {
			buf, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if globalMap.GetActiveChannel() == client.id {
					fmt.Printf("Client %d disconnected. Removing from list...\n", client.id)
				}
				// We need the kill flag here because sometimes
				// runReadClient executes before the KillClient()
				// functtion can complete. So, we use killFlag
				// to make sure that isn't currently happening.
				if !globalMap.IsDead(client.id) && !killFlag {
					clientDisconnect(client)
				}
				caret()
				return
			}
			buf = b64_decode(buf)
			fmt.Print(buf)
			caret()
		} else {
			buf := make([]byte, 65535)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Client %d disconnected. Removing from list...\n", client.id)
				caret()
				if !globalMap.IsDead(client.id) {
					clientDisconnect(client)
				}
				return
			}
			if globalMap.GetActiveChannel() == client.id {
				fmt.Print(string(buf))
			}
		}
	}
}

func runWriteClient(client Client) {
	conn := client.conn
	channel := globalMap.GetChannel(client.id)
	isActive := <-*channel
	for {
		if isActive && globalMap.GetActiveChannel() == client.id {
			if globalMap.IsDead(client.id) {
				return
			}
			stdin := globalMap.GetStdin()
			stdReadLine := <-*stdin
			if stdReadLine == "!!!FIN!!!" {
				isActive = false
				// return
			} else if stdReadLine == "!!!DEAD!!!" {
				return
			} else {
				conn.Write([]byte(fmt.Sprintf("%s\n", trim(stdReadLine))))
			}
		} else {
			if globalMap.IsDead(client.id) {
				return
			}
			isActive = <-*channel
		}
	}

}
