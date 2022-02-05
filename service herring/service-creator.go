//Michael Burke, mdb5315@rit.edu
package main

import  (
	"ioutil"
	"fmt"
	)


func main() {
	dat, _ := ioutil.ReadFile("template.service")
	file := string(dat)
	fmt.Printf(file)
}