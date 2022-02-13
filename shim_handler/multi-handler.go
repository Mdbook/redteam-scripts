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
	wg.Wait()
	fmt.Println("Finished!")
}

func do(ip string) {
	defer wg.Done()
	listenPort := strings.ReplaceAll(ip, ".", "")
	listenPort = "2" + listenPort[len(listenPort)-4:]
	cmd := exec.Command("xterm", "-title", ip, "-e", "nc", "-l", "-p", listenPort)
	cmd.Run()
}
