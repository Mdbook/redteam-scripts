package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	resp, err := http.Get("http://mdbook.me/ip.txt")
	var ip string
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		line := string(body)
		line = strings.TrimSuffix(line, "\n")
		ip = line
	} else {
		resp, err = http.Get("http://129.21.141.218/ip.txt")
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			line := string(body)
			line = strings.TrimSuffix(line, "\n")
			ip = line
		} else {
			ip = "10.0.100.101"
		}
	}
	fmt.Print(ip)
}
