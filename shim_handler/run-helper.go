package main

import (
	"fmt"
	"os/exec"
	"sync"
)

var wg sync.WaitGroup

func openXTerm(title, file string) {
	defer wg.Done()
	for {
		cmd := exec.Command("xterm", "-title", title, "-e", "go run "+file)
		cmd.Run()
	}
}

func main() {
	fmt.Println("Starting handlers...")
	wg.Add(4)
	go openXTerm("LS | MASTER", "ls-server.go")
	go openXTerm("VI | MASTER", "vi-server.go")
	go openXTerm("VIM | MASTER", "vim-server.go")
	go openXTerm("NANO | MASTER", "nano-server.go")
	fmt.Println("Handlers started")
	wg.Wait()
}
