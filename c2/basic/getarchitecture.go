package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Your OS is:	         %s\n", runtime.GOOS)
	fmt.Printf("Your architecture is:    %s\n", runtime.GOARCH)
}
