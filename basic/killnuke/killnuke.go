//Michael Burke
//Simple program to render it more difficult to kill jumper.

package main

import (
	"os"
	"os/exec"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-o", "-y":
				return
			}
		}
	}
	cmd := exec.Command("ki11a11", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
