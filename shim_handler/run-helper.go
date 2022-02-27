package main

import (
	"fmt"
	"os/exec"
	"sync"
)

var wg sync.WaitGroup

func openXTerm(title, file string) {
	// wg.Add(1)
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

// xterm -title "LS | MASTER" -e "go run ls-server.go" | &
// xterm -title "VI | MASTER" -e "go run vi-server.go" | &
// xterm -title "VIM | MASTER" -e "go run vim-server.go" | &
// xterm -title "NANO | MASTER" -e "go run nano-server.go" | &
