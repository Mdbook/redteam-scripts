package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// var port string
var hasPort bool = false
var takenPorts []string
var wg sync.WaitGroup
var globalMap globalMaster
var stdin chan string
var HOST_IP string

func main() {
	handleArgs()
	go readStdin()
	// return
	// TODO: add handler for multiple ports
	fmt.Println("Server starting...")
	for {
		connectionHelper()
	}
}

func handleArgs() {
	HOST_IP = GetOutboundIP()
	globalMap = *CreateMaster()
	// args := os.Args
	// if len(args) > 1 {
	// 	for i := 1; i < len(args); i++ {
	// 		switch args[i] {
	// 		case "-p":
	// 			port = args[i+1]
	// 			hasPort = true
	// 		}
	// 	}
	// }
	// if !hasPort {
	// 	fmt.Printf("Port: ")
	// 	fmt.Scanln(&port)
	// 	hasPort = true
	// }
}

func connectionHelper() {
	wg.Add(4)
	go GetConnection("8003")
	go GetConnection("8004")
	go GetConnection("8005")
	go GetConnection("8006")
	wg.Wait()

}

func GetConnection(port string) {
	defer wg.Done()
	for {
		getPort, _ := net.Listen("tcp", HOST_IP+":"+port)

		conn, _ := getPort.Accept()
		var clientInfo ClientInfo
		remoteInfo, err := bufio.NewReader(conn).ReadString('\n')
		if strings.Index(remoteInfo, "INFO:") != -1 {
			fmt.Println(remoteInfo)
			clientInfo = parseParams(remoteInfo)
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		remotePort := getRandomPort()
		takenPorts = append(takenPorts, remotePort)

		go handleClient(clientInfo, remotePort)
		time.Sleep(100 * time.Millisecond)
		conn.Write([]byte(remotePort))
		conn.Close()
		getPort.Close()
	}
}

func enterTerminal(channel *chan string, reader *bufio.Reader) {
	fmt.Printf("Entered terminal for client %d. Type 'leave' to leave.\n", globalMap.GetActiveChannel())
	for {
		cmd, _ := reader.ReadString('\n')
		if trim(cmd) == "leave" {
			fmt.Println("---Leaving terminal---")
			return
		}
		*channel <- cmd
	}

}

func readStdin() {
	channel := globalMap.GetStdin()
	for {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = trim(cmd)
		args := strings.Split(cmd, " ")
		switch args[0] {
		case "exit":
			os.Exit(0)
		case "send":
			globalMap.SetSingle(true)
			*channel <- cmd[5:]
		case "enter":
			globalMap.SetSingle(false)
			fmt.Println("---Entering terminal---")
			enterTerminal(channel, reader)
		case "set":
			switch args[1] {
			case "active":
				id, _ := strconv.Atoi(args[2])
				if globalMap.GetActiveChannel() == id {
					fmt.Printf("Channel is already active!\n\n")
					break
				} else if globalMap.GetCurrentId() <= id || id < 0 {
					fmt.Printf("Error: index out of range\n\n")
					break
				}
				if globalMap.GetActiveChannel() != -1 {
					*channel <- "!!!FIN!!!"
				}
				globalMap.SetActive(id)
				fmt.Printf("Set client %d as active\n", id)
				fmt.Println()
			default:
				displayHelp("set")
			}
		case "get":
			switch args[1] {
			case "client":
				var clientId int
				if len(args) < 3 {
					clientId = globalMap.GetActiveChannel()
				} else {
					clientId, _ = strconv.Atoi(args[2])
				}
				// fmt.Println(globalMap.GetCurrentId())
				// fmt.Println(clientId)
				if globalMap.GetCurrentId() <= clientId || clientId < 0 {
					if len(args) < 3 {
						fmt.Println("Error: no client active")
					} else {
						fmt.Println("Error: invalid client ID")
					}
				} else {
					fmt.Printf("Info for client %d:\n", clientId)
					curClient := globalMap.GetClient(clientId)
					fmt.Printf(
						"LAN IP: %s\n"+
							"WAN IP: %s\n"+
							"Port: %s\n"+
							"Using client: %s\n"+
							"OS: %s\n"+
							"OS Type: %s\n"+
							"Encoded connection: %t\n",
						curClient.lanIP,
						curClient.wanIP,
						curClient.port,
						curClient.clientType,
						curClient.os,
						curClient.osFlavor,
						curClient.isEncoded,
					)
				}
				fmt.Println()
			case "clients":
				fmt.Println("Current clients: ")
				for _, client := range globalMap.GetClients() {
					fmt.Printf("Client %d\n", client.id)
				}
				fmt.Println()
			case "active":
				switch args[2] {
				case "client":
					fmt.Printf("Current active client: %d\n", globalMap.GetActiveChannel())
				}
				fmt.Println()
			default:
				displayHelp("get")
			}
		case "help":
			displayHelp(cmd)
		default:
			fmt.Println("Error: unknown command")
		}

		// if cmd[0] == '>' {
		// 	isKill, id := handleCommand(cmd[1:])
		// 	if isKill {

		// 	}
		// } else {
		//
		// }
	}
}

func handleCommand(cmd string) (bool, int) {
	return false, 0
}

func do(client Client) {
	channel := globalMap.GetChannel(client.id)
	isActive := <-*channel
	for {
		if isActive && globalMap.GetActiveChannel() == client.id {
			stdin := globalMap.GetStdin()
			stdReadLine := <-*stdin
			if stdReadLine == "!!!FIN!!!" {
				isActive = false
			} else {
				fmt.Printf("Received on client %d: %s\n", client.id, stdReadLine)
			}
		} else {
			isActive = <-*channel
		}
	}

}

func displayHelp(cmd string) {

}

func getRandomPort() string {
	port1 := strconv.Itoa(random(10))
	port2 := strconv.Itoa(random(99))
	if len(port2) <= 1 {
		port2 = "0" + port2
	}
	// fmt.Println(port1)
	// fmt.Println(port2)
	remotePort := "2" + "5" + port1 + port2
	if findIndex(takenPorts, remotePort) == -1 {
		return remotePort
	}
	return getRandomPort()
}
