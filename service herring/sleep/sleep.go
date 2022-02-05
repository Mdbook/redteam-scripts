package main

import "time"

func main() {
	do()
}

func do() {
	time.Sleep(2000 * time.Second)
	do()
}
