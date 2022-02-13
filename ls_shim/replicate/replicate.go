package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

var systemOS string = getOS()

func main() {
	fmt.Println(systemOS)
	//if systemOS == "centos"
}

func getOS() string {
	var ret_os string
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
		if strings.Index(os_split[i], "ID_LIKE=") != -1 {
			ret_os = strings.Replace(os_split[i], "ID=", "", 1)
			ret_os = strings.Replace(ret_os, `"`, "", -1)
			break
		}
	}
	return ret_os
}

func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
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
