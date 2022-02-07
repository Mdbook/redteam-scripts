package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

var messages []string

func main() {
	buildMessages()
	do()
}

func do() {
	rand.Seed(time.Now().UnixNano())
	delay := (random(19) + 1) * 60
	fmt.Printf("Sleeping for %d Minutes\n", delay/60)
	time.Sleep(time.Duration(delay) * time.Second)
	fmt.Printf("Sending Message\n")
	index := random(len(messages) - 1)
	command := messages[index]
	cmd := exec.Command("wall", command)
	cmd.Run()
	do()
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func buildMessages() {
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
		"Can you hurry up and fine me already? I'm bored",
		"HehehehehehehehEHEHEHEHEHEHE",
	}
}
