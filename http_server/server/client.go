package main

import (
	"fmt"
	"net"
)

type Client struct {
	id     int
	ip     string
	port   string
	client string
	conn   net.Conn
}

func handleClient(port, remoteClient string) {
	getPort, _ := net.Listen("tcp", GetOutboundIP()+":"+port)
	conn, _ := getPort.Accept()
	// defer conn.Close()
	// defer getPort.Close()
	remoteIp := conn.RemoteAddr().String()
	// fmt.Printf("Received request from %s\n", remoteIp)
	client := globalMap.CreateClient(remoteIp, port, remoteClient, conn)
	fmt.Printf("New client (%s) connected with id: %d \n", remoteIp, client.id)
	go runReadClient(client)
	go runWriteClient(client)
}

func runReadClient(client Client) {
	conn := client.conn
	conn.Write([]byte(fmt.Sprintf("%s\n", "dir")))
	for {
		buf := make([]byte, 65535)
		_, err := conn.Read(buf)
		if err != nil {
			// TODO error handling
			fmt.Println(err)
			break
		}
		if globalMap.GetActiveChannel() == client.id {
			fmt.Print(string(buf))
			if globalMap.IsSingle() {
				fmt.Println()
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
			stdin := globalMap.GetStdin()
			stdReadLine := <-*stdin
			if stdReadLine == "!!!FIN!!!" {
				isActive = false
			} else {
				conn.Write([]byte(fmt.Sprintf("%s\n", stdReadLine)))
			}
		} else {
			isActive = <-*channel
		}
	}

}

func readConnection(conn net.Conn, active *bool, channel chan string) {
	defer wg.Done()
	for {
		buf := make([]byte, 65535)
		_, err := conn.Read(buf)
		if err != nil {
			channel <- "!!!FIN!!!"
			*active = false
			fmt.Println(err)
			break
		}
		fmt.Print(string(buf))
	}
}

func writeConnection(conn net.Conn, active *bool, channel chan string) {
	defer wg.Done()
	for {
		if *active {
			text := <-channel
			if text == "!!!FIN!!!" {
				return
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", text)))
		} else {
			return
		}

	}

}
