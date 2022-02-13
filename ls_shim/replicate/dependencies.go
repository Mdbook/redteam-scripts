package main

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

var systemOS string = getOS()

func main() {
	installDependencies()
}

func getOS(isFail ...bool) string {
	var ret_os string
	checkID := false
	if len(isFail) > 0 && isFail[0] {
		checkID = true
	}
	os_str := readFile("/etc/os-release")
	os_split := strings.Split(os_str, "\n")
	for i := 0; i < len(os_split); i++ {
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
		return getOS(true)
	}
	return ret_os
}

func installDependencies() {
	if systemOS == "debian" {
		cmd := exec.Command("apt-get", "install", "golang-go", "-y")
		cmd.Run()
	} else if systemOS == "arch" {
		cmd := exec.Command("pacman", "-S", "go", "--noconfirm")
		cmd.Run()
	} else if strings.Index(systemOS, "rhel") != -1 {
		cmd := exec.Command("yum", "install", "golang", "-y")
		cmd.Run()
	} else if systemOS == "fedora" {
		cmd := exec.Command("dnf", "install", "golang", "-y")
		cmd.Run()
	}
}

func readFile(path string) string {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return str
}
