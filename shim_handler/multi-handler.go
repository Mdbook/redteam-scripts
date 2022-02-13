package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	args := os.Args
	for i := 1; i < len(args); i++ {
		ip := args[i]
		wg.Add(1)
		go do(ip)
	}
	wg.Add(1)
	go mainServer()
	wg.Wait()
	fmt.Println("Finished!")
}

func do(ip string) {
	defer wg.Done()
	listenPort := strings.ReplaceAll(ip, ".", "")
	listenPort = "2" + listenPort[len(listenPort)-4:]
	for {
		cmd := exec.Command("xterm", "-title", ip, "-e", "nc", "-l", "-p", listenPort)
		cmd.Run()
	}

}

func mainServer() {
	defer wg.Done()
	for {
		cmd := exec.Command("xterm", "-title", "master", "-e", "go", "run", "network-server.go")
		cmd.Run()
	}

}
