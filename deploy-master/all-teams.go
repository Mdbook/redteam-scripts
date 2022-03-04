package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func runDeployWAN(i, password string) {
	defer wg.Done()
	fmt.Println("Deploying to 172.16." + i + ".0/24 (DETACHED)")
	cmd := exec.Command("xterm", "-e", "go", "run", "deploy-master.go", "-v", "-t", "172.16."+i+".20", "--user-list", "karen,joyce,raymond,howard,suzanne,louise,glenn,christopher,lynn,enzo,peggy,margaret,melvin,wendell", "-p", password)
	cmd.Run()
}

func runDeployLAN(i, password string) {
	defer wg.Done()
	fmt.Println("Deploying to 10." + i + ".1.0/24 (DETACHED)")
	cmd := exec.Command("xterm", "-e", "go", "run", "deploy-master.go", "-v", "-t", "10."+i+".1.69", "--user-list", "karen,joyce,raymond,howard,suzanne,louise,glenn,christopher,lynn,enzo,peggy,margaret,melvin,wendell", "-p", password)
	cmd.Run()
}

func main() {
	numTeams := 15
	var curTeam int
	var password string
	fmt.Printf("Your team: ")
	_, _ = fmt.Scanf("%d", &curTeam)
	fmt.Println()
	fmt.Printf("Password: ")
	_, _ = fmt.Scanf("%s", &password)
	fmt.Println()
	for i := 1; i <= numTeams; i++ {
		if i != curTeam {
			wg.Add(1)
			go runDeployLAN(strconv.Itoa(i), password)
		}
	}
	for i := 1; i <= numTeams; i++ {
		if i != curTeam {
			wg.Add(1)
			go runDeployWAN(strconv.Itoa(i), password)
		}
	}
	wg.Wait()
	fmt.Println("Finished deploying to all systems.")
}
