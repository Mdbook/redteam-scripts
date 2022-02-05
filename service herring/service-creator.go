//Michael Burke, mdb5315@rit.edu
package main

import  (
	"io/ioutil"
	"fmt"
	)


func createServices(){
	names := {"yourmom", "grap", "amogus", "sus", "virus", "redteam", "the-matrix", "uno-reverse-card", "yellowteam", "bingus", "dokidoki", "based", "not-ransomware", "bepis", "roblox", "freevbucks", "notavirus", "heckerman", "benignfile", "yolo", "pickle", "grubhub", "hehe", "amogOS", "society", "yeet", "doge", "mac", "hungy", "youllneverfindme", "red-herring"}
	descriptions = {
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
		"Hot singles in your area"
	}
	

}


func main(){
	dat, _ := ioutil.ReadFile("template.service")
	file := string(dat)
	fmt.Printf(file)
}