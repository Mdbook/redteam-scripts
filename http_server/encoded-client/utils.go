package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type ClientInfo struct {
	lanIP      string
	clientType string
	os         string
	osFlavor   string
	isEncoded  bool
}

func parseParams(info string) ClientInfo {
	info = info[strings.Index(info, "INFO:{")+6 : strings.Index(info, "}")]
	params := strings.Split(info, ",")
	clientInfo := ClientInfo{}
	fmt.Println(clientInfo.clientType)
	for _, paramString := range params {
		param := strings.Split(paramString, ":")
		switch param[0] {
		case "isEncoded":
			clientInfo.isEncoded, _ = strconv.ParseBool(param[1])
		case "clientType":
			clientInfo.clientType = param[1]
		case "lanIP":
			clientInfo.lanIP = param[1]
		case "os":
			clientInfo.os = param[1]
		case "osFlavor":
			clientInfo.osFlavor = param[1]
		}
	}
	return clientInfo
}

func debugf(s string, params ...interface{}) {
	fmt.Printf(s, params...)
}

func debugln(s interface{}) {
	fmt.Println(s)
}

func trim(str string) string {
	return strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
}

func getOS(isFail ...bool) string {
	var ret_os string
	checkID := false
	//Nifty little code to let us call getOS with or without a parameter
	if len(isFail) > 0 && isFail[0] {
		checkID = true
	}
	//Read /etc/os-release to find what distro the host is running on
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
		//Child distros will have ID_LIKE instead of ID. Check for both
		matchString := "ID_LIKE="
		if checkID {
			matchString = "ID="
		}
		if strings.Index(os_split[i], matchString) == 0 {
			ret_os = strings.Replace(os_split[i], matchString, "", 1)
			ret_os = strings.Replace(ret_os, `"`, "", -1)
			break
		}
	}
	if ret_os == "" && !checkID {
		// If ID_LIKE wasn't found, then seach for ID= instead
		return getOS(true)
	}
	return ret_os
}

func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
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

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
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
