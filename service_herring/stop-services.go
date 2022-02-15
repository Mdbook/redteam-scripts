package main

import (
	"fmt"
	"os/exec"
)

func main() {
	names := []string{"yourmom", "freddy-fazbear", "grap", "amogus", "sus", "virus", "redteam", "the-matrix", "uno-reverse-card", "yellowteam", "bingus", "dokidoki", "based", "not-ransomware", "bepis", "roblox", "freevbucks", "notavirus", "heckerman", "benignfile", "yolo", "pickle", "grubhub", "hehe", "amogOS", "society", "yeet", "doge", "mac", "hungy", "youllneverfindme", "red-herring"}
	for i := 0; i < len(names); i++ {
		cmd := exec.Command("systemctl", "stop", names[i]+".service")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Couldn't stop " + names[i] + ".service")
		} else {
			fmt.Println("Stopped " + names[i] + ".service")
		}

	}
	fmt.Println("All services stopped.")
}
