package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type Colors struct {
	red     string
	green   string
	blue    string
	black   string
	yellow  string
	magenta string
	cyan    string
	white   string
	reset   string
}

func initColors() Colors {
	return Colors{
		reset:   "\033[0m",
		black:   "\033[30m",
		red:     "\033[31m",
		green:   "\033[32m",
		yellow:  "\033[33m",
		blue:    "\033[34m",
		magenta: "\033[35m",
		cyan:    "\033[36m",
		white:   "\033[37m",
	}
}
func parseParams(info string) ClientInfo {
	info = info[strings.Index(info, "INFO:{")+6 : strings.Index(info, "}")]
	info = b64_decode(info)
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
		case "hostname":
			clientInfo.hostname = param[1]
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

func handleQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			sig.Signal()
			fmt.Println("\nType \"exit\" to exit")
			caret()
			// sig is a ^C, handle it
		}
	}()
}

func errorln(str string) {
	fmt.Println(colors.red + str + colors.reset)
}

func infoln(str string) {
	fmt.Println(colors.yellow + str + colors.reset)
}

func caret() {
	fmt.Print(colors.green + "> " + colors.reset)
}
