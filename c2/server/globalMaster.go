package main

import (
	"net"
	"sync"
)

type globalMaster struct {
	clients       []Client
	channels      []chan bool
	stdin         chan string
	activeChannel int
	currentId     int
	isSingle      bool
	mux           sync.Mutex
}

func CreateMaster() *globalMaster {
	return &globalMaster{activeChannel: -1, currentId: 0, isSingle: false, stdin: make(chan string)}
}

func (a *globalMaster) IsSingle() bool {
	a.mux.Lock()
	isSingle := a.isSingle
	a.mux.Unlock()
	return isSingle
}

func (a *globalMaster) SetSingle(isSingle bool) {
	a.mux.Lock()
	a.isSingle = isSingle
	a.mux.Unlock()
}

func (a *globalMaster) GetClients() []Client {
	a.mux.Lock()
	clients := a.clients
	a.mux.Unlock()
	return clients
}

func (a *globalMaster) GetClient(id int) Client {
	a.mux.Lock()
	clients := a.clients[id]
	a.mux.Unlock()
	return clients
}

func (a *globalMaster) GetChannel(id int) *chan bool {
	a.mux.Lock()
	channel := &a.channels[id]
	a.mux.Unlock()
	return channel
}

func (a *globalMaster) GetActiveChannel() int {
	a.mux.Lock()
	id := a.activeChannel
	a.mux.Unlock()
	return id
}

func (a *globalMaster) GetStdin() *chan string {
	a.mux.Lock()
	stdin := &a.stdin
	a.mux.Unlock()
	return stdin
}

func (a *globalMaster) SetActive(id int) {
	a.mux.Lock()
	// if a.activeChannel != -1 {
	// a.stdin <- ""
	// a.channels[a.activeChannel] <- false
	// }
	a.channels[id] <- true
	a.activeChannel = id
	a.mux.Unlock()
}

func (a *globalMaster) CreateClient(clientInfo ClientInfo, port string, conn net.Conn, listener net.Listener) Client {
	a.mux.Lock()
	cliId := a.currentId
	a.channels = append(a.channels, make(chan bool))
	a.currentId++
	client := Client{
		id:         cliId,
		lanIP:      clientInfo.lanIP,
		wanIP:      conn.RemoteAddr().String(),
		port:       port,
		clientType: clientInfo.clientType,
		os:         clientInfo.os,
		osFlavor:   clientInfo.osFlavor,
		hostname:   clientInfo.hostname,
		isEncoded:  clientInfo.isEncoded,
		isDead:     false,
		conn:       conn,
		listener:   listener,
	}
	a.clients = append(a.clients, client)
	a.mux.Unlock()
	return client
}

func (a *globalMaster) GetCurrentId() int {
	a.mux.Lock()
	id := a.currentId
	a.mux.Unlock()
	return id
}
func (a *globalMaster) SetCurrentId(id int) {
	a.mux.Lock()
	a.currentId = id
	a.mux.Unlock()
}

func (a *globalMaster) IsDead(id int) bool {
	a.mux.Lock()
	dead := a.clients[id].isDead
	a.mux.Unlock()
	return dead
}

func (a *globalMaster) KillClient(id int) {
	a.mux.Lock()
	a.clients[id].isDead = true
	if a.activeChannel == id {
		a.activeChannel = -1
		a.stdin <- "!!!FIN!!!"
		// a.channels[id] <- false
	} else {
		a.channels[id] <- true
	}
	a.mux.Unlock()
}
