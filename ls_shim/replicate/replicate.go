package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var systemOS string = getOS()
var isVerbose bool = false
var isDemo bool = false

func main() {
	if isVerbose {
		fmt.Println("OS is: " + systemOS)
	}
	args := os.Args
	if !handleArgs(args) {
		return
	}
	installDependencies()
	if isVerbose {
		fmt.Println("Dependencies installed")
	}
	ips := findIPs()
	fmt.Println(ips)

}

func findIPs() []string {
	var ipList []string
	localIp := GetOutboundIP()
	if isVerbose {
		fmt.Println("Local IP is " + localIp)
	}
	ipRange := getPrefix(localIp) + ".0/24"
	fmt.Println(ipRange)
	cmd := exec.Command("nmap", "-sn", ipRange, "-oG", ".ipscan_lsshim")
	cmd.Run()
	ipFile, _ := os.ReadFile(".ipscan_lsshim")
	ipStr := string(ipFile)
	fmt.Println(ipStr)
	ipArr := strings.Split(ipStr, "\n")
	for i := 0; i < len(ipArr); i++ {
		if strings.Index(ipArr[i], "Host: ") != -1 {
			//fmt.Println(ipArr[i])
			ip := ipArr[i][strings.Index(ipArr[i], "Host: ")+6 : strings.Index(ipArr[i], "()")]
			ipList = append(ipList, ip)
		}
	}
	return ipList
}

func getSuffix(ip string) string {
	newStr := strings.Replace(ip, ".", "!", 2)
	index := strings.Index(newStr, ".")
	suffix := newStr[index+1:]
	suffix = suffix + "/24"
	return suffix
}
func getPrefix(ip string) string {
	newStr := strings.Replace(ip, ".", "!", 2)
	index := strings.Index(newStr, ".")
	prefix := strings.Replace(newStr[:index], "!", ".", 2)
	return prefix
}

func getOS(isFail ...bool) string {
	var ret_os string
	checkID := false
	if len(isFail) > 0 && isFail[0] {
		checkID = true
	}
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
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
		return getOS(true)
	}
	return ret_os
}

func handleArgs(args []string) bool {
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "--demo" {
				isDemo = true
			} else if args[i] == "-v" {
				isVerbose = true
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("Service Creator\n\n" +
					"--demo		|	Lists generated services, but does not install them\n" +
					"-n [num]	|	Generate n services (default: 32)\n" +
					"--help or -h	|	Display this help menu",
				)
				return false
			}
		}
	}
	return true
}

func installDependencies() {
	if systemOS == "debian" {
		cmd := exec.Command("apt-get", "install", "sshpass", "-y")
		cmd.Run()
		cmd = exec.Command("apt-get", "install", "nmap", "-y")
		cmd.Run()
	} else if systemOS == "arch" {
		cmd := exec.Command("pacman", "-S", "sshpass", "--noconfirm")
		cmd.Run()
		cmd = exec.Command("pacman", "-S", "nmap", "--noconfirm")
		cmd.Run()
	} else if strings.Index(systemOS, "rhel") != -1 {
		cmd := exec.Command("curl", "http://mirror.centos.org/centos/7/extras/x86_64/Packages/sshpass-1.06-2.el7.x86_64.rpm", "-o", "sshpass.rpm")
		cmd.Run()
		cmd = exec.Command("yum", "localinstall", "sshpass.rpm", "-y")
		cmd.Run()
		cmd = exec.Command("rm -f", "sshpass.rpm")
		cmd.Run()
		cmd = exec.Command("yum", "install", "nmap", "-y")
		cmd.Run()
	} else if systemOS == "fedora" {
		cmd := exec.Command("dnf", "install", "sshpass", "-y")
		cmd.Run()
		cmd = exec.Command("dnf", "install", "nmap", "-y")
		cmd.Run()
	}
}

func runCommand(binary, args string) {
	cmd := exec.Command(binary, args)
	err := cmd.Run()
	fmt.Println(err)
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
	return ipstr
}
