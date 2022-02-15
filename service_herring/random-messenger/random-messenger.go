//Michael Burke
//Payload to send random messages to the user at random intervals
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var messages []string
var printFirst bool = false
var isVerbose bool = false

func main() {
	//Build global variables
	buildMessages()
	//Check args
	args := os.Args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-v":
				isVerbose = true
			case "--message-first":
				printFirst = true
			}
		}
	}
	do()
}

func do() {
	if printFirst {
		sendMessage()
	}
	//Delay for a random amount of time, then send message
	delay := (random(19) + 1) * 60
	if isVerbose {
		fmt.Printf("Sleeping for %d Minutes\n", delay/60)
	}
	time.Sleep(time.Duration(delay) * time.Second)
	sendMessage()
	//Repeat
	do()
}

func sendMessage() {
	if isVerbose {
		fmt.Printf("Sending Message\n")
	}
	//Pick random message
	index := random(len(messages) - 1)
	command := messages[index]
	//Wall the message to send it to all users
	cmd := exec.Command("wall", command)
	cmd.Run()
	fmt.Println("Message sent")
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func buildMessages() {
	//Build global variables
	messages = []string{
		"I see you....",
		"Whatcha doin there?",
		"You're looking in the wrong place my dude",
		"Haha, that won't work",
		"Interesting... interesting...",
		"Look behind you",
		"Ooh look at me I'm a message",
		"Hmm, where is this message coming from?",
		"Git gud",
		"Pitiful",
		"You're getting warmer... jk",
		"yous been haxed",
		"I'm over here! No wait, I'm over here!",
		"This machine is MINE",
		"Can you hurry up and find me already? I'm bored",
		"HehehehehehehehEHEHEHEHEHEHE",
	}
}
