//Michael Burke
//Simple payload that does literally nothing

package main

import "time"

func main() {
	do()
}

func do() {
	//Sleep for 2000 seconds
	time.Sleep(2000 * time.Second)
	do()
}
