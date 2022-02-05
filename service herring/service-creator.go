//Michael Burke, mdb5315@rit.edu
package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
)

type service struct {
	name        string
	description string
	path        string
	filename    string
	payload     string
	user        string
}

var user string
var names, descriptions, paths, filenames, payloads []string

func buildDB() {
	user = "root"
	names = []string{"yourmom", "freddy-fazbear", "grap", "amogus", "sus", "virus", "redteam", "the-matrix", "uno-reverse-card", "yellowteam", "bingus", "dokidoki", "based", "not-ransomware", "bepis", "roblox", "freevbucks", "notavirus", "heckerman", "benignfile", "yolo", "pickle", "grubhub", "hehe", "amogOS", "society", "yeet", "doge", "mac", "hungy", "youllneverfindme", "red-herring"}
	descriptions = []string{
		"An absolutely vital service for Linux. Do not delete under any circumstances. -redteam",
		"kinda sus bro",
		"Very benign. Much trust.",
		"uhhhhhhh",
		"malware go brrrr",
		"Smudge the crunchy cat",
		"Do me a favor and keep this service running, I have a wife and 4 kids to feed",
		"We've been trying to reach you about your car's extended warranty",
		"hehehehehehehehehehehe",
		"UwU what's this?",
		"Vanessa, I'm a material gorl!",
		"I turned myself into a service morty! I'm service Rick!",
		"If you or a love one has been diagnosed with mesothelioma, you may be entitled to a cash reward",
		"It's free real estate",
		"Hot singles in your area",
		"Meesa jar jar binks",
	}
	paths = []string{
		"/var/run/",
		"/var/",
		"/etc/",
		"/home/",
		"/usr/lib/",
		"/usr/local/",
		"/root/",
	}
	filenames = []string{
		"randomservice",
		"inconspicuous_file",
		"deleteme",
		"dontdeleteme",
		"heh",
		"b1ngus",
		"file12345",
		"homework",
		"top-secret",
		"temporary-file",
		"lilboi",
		"geck",
		"flappy-bird-game",
		"borger",
		"issaservice",
		"himom",
		"jeffUwU",
		"youfoundme",
	}
	payloads = []string{"random-messenger", "reverse-shell", "downloader", "file-creator", "user-creator"}
}

func buildServices(num int) []service {
	validNames := names
	var services []service
	for i := 0; i < num; i++ {
		var serviceName, serviceDesc, servicePath, serviceFilename, servicePayload string
		validNames, serviceName = pickFrom(validNames)
		serviceDesc = getRandom(descriptions)
		servicePath = getRandom(paths)
		serviceFilename = getRandom(filenames)
		servicePayload = getRandom(payloads)
		for {
			if hasCollision(services, servicePath, serviceFilename) {
				servicePath = getRandom(paths)
				serviceFilename = getRandom(filenames)
			} else {
				break
			}
		}
		newService := service{serviceName, serviceDesc, servicePath, serviceFilename, servicePayload, user}
		services = append(services, newService)
	}
	return services
}

func hasCollision(services []service, servicePath string, serviceFilename string) bool {
	for i := 0; i < len(services); i++ {
		curService := services[i]
		if curService.path == servicePath && curService.filename == serviceFilename {
			return true
		}
	}
	return false
}

func pickFrom(slice []string) ([]string, string) {
	var val string
	slice, val = remove(slice, getRandomIndex(slice))
	return slice, val
}

func getRandomIndex(slice []string) int {
	return rand.Intn(len(slice) - 1)
}

func getRandom(slice []string) string {
	randNum := rand.Intn(len(slice) - 1)
	return slice[randNum]
}

func remove(slice []string, i int) ([]string, string) {
	name := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, name

}

func findIndex(slice []string, value string) int {
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

func main() {
	buildDB()
	dat, _ := ioutil.ReadFile("template.service")
	file := string(dat)
	buildServices(len(names))
	fmt.Println(file)
	fmt.Println(len(names))
}
