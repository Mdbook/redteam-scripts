//Michael Burke, mdb5315@rit.edu
package main

import  (
	"os/exec"
	"os"
	)
func main() {
	rev := exec.Command("systemd-restart")
	rev.Run()
	cmd := exec.Command("lsa", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}