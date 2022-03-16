package main

import (
	"fmt"
	"sync"
	"time"
)

type uni struct {
	channels sync.Map
}

var globalMap uni = uni{}
var test chan string = make(chan string)
var wg sync.WaitGroup
var globalChannelId int = 9

func main() {
	makeChannels()
}

func makeChannels() {
	defer close(test)
	for i := 0; i < 10; i++ {
		channel := make(chan bool)
		wg.Add(1)
		go storeChannel(i, channel)
	}
	wg.Wait()
	for i := 0; i < 10; i++ {
		go do(i)
	}
	a, _ := globalMap.channels.Load(globalChannelId)
	channel := a.(chan bool)
	channel <- true
	time.Sleep(100 * time.Millisecond)
}
func storeChannel(i int, channel chan bool) {
	defer wg.Done()
	globalMap.channels.Store(i, channel)
}

func do(id int) {
	a, _ := globalMap.channels.Load(id)
	curChannel := a.(chan bool)
	_ = <-curChannel
	fmt.Println(id)
}
