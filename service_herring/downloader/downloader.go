package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var downDir string = "./"
var host string = "http://192.168.12.6/"

func main() {
	do()
}

func do() {
	stat := getHTTP(host + "stat")
	commands := strings.Split(stat, "\n")
	for i := 0; i < len(commands); i++ {
		fmt.Println("Checking " + commands[i])
		command := commands[i]
		switch command {
		case "DOWNLOAD":
			fmt.Println("Downloading!")
			downloadFile(downDir, host+commands[i+1])
		}

	}
	//repeat()
}

func repeat() {
	delay := (random(19) + 1) * 60
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
