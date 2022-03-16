package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type ClientInfo struct {
	lanIP      string
	clientType string
	isEncoded  bool
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

func remove(slice []string, i int) ([]string, string) {
	//Remove an item from a slice
	name := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, name
}

func findIndex(slice []string, value string) int {
	//Find the index of a string in a slice
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

func parseParams(info string) ClientInfo {
	info = info[strings.Index(info, "INFO:{")+6 : strings.Index(info, "}")]
	params := strings.Split(info, ",")
	clientInfo := ClientInfo{}
	fmt.Println(clientInfo.clientType)
	for _, paramString := range params {
		param := strings.Split(paramString, ":")
		switch param[0] {
		case "clientType":
			clientInfo.clientType = param[1]
		case "lanIP":
			clientInfo.lanIP = param[1]
		case "isEncoded":
			clientInfo.isEncoded, _ = strconv.ParseBool(param[1])
		}
	}
	return clientInfo
}

/**
Base 64 encode a message to be sent to the server
*/
func b64_encode(text string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	return encoded
}

/**
base 64 decode a message from the server
*/
func b64_decode(text string) string {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		fmt.Println(err)
	}
	return string(decoded)
}
