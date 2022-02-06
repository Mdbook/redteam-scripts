package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var users []string
var password string
var paths []string

func main() {
	buildDB()
	//fmt.Println(getRecursive("/etc/"))
	do()

}

func do() {
	filename := randString(random(19)+1) + "." + randString(random(4)+1)
	path := getPath()
	fmt.Println(path + filename)
	//writeToFile(path, filename)

}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func getRecursive(path string) string {
	//fmt.Println(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var directories []string
	for _, f := range files {
		if f.IsDir() && f.Name()[:1] != "." {
			directories = append(directories, f.Name())
		}
	}
	if len(directories) == 0 {
		return path
	}
	for i := 0; i < len(directories); i++ {
		choose := random(len(directories) + 1)
		if choose == 1 || choose == len(directories) {
			return getRecursive(path + directories[i] + "/")
		}
	}
	return path
}

func getPath() string {
	path := getRandom(paths)
	if strings.Index(path, "*") != -1 {
		return getRecursive(strings.ReplaceAll(path, "*", ""))
	}
	return path
}

func getRandom(slice []string) string {
	if len(slice) == 1 {
		return slice[0]
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(slice) - 1)
	return slice[randNum]
}

func writeToFile(path string, str string) {
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString("Hello World")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func buildDB() {
	paths = []string{
		"/etc/*",
		"/home/*",
		"/mnt/",
		"/root/",
		"/var/run/",
		"/usr/lib",
		"/usr/bin/",
		"/etc/",
		"/home/",
		"/usr/*",
		"/var/log/*",
		"/var/log/",
	}
}
