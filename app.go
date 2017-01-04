package main

import (
	"os"
	"fmt"
	"./api"
	"net"
	"time"
	"math/rand"
)

// Registering callbacks for the state from across every files
func registerAllCallbacks(self *api.State) {
	api.RegisterLeaderCallbacks(self)
}


func processLogic(self *api.State) {
	registerAllCallbacks(self)
	go api.Listen(self)
	go api.CheckLeader(self)

	// This will be called after a leader has been elected
	//go api.GenerateInput(self)

	//go api.PingEveryone(self)
	//call the other functions as they are made
}

func main() {
	if len(os.Args) < 3 {
		println("Wrong usage")
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())
	// The state was changed to a pointer to allow its manipulation within functions
	var self *api.State = new(api.State)
	self.ListenPort = os.Args[1]
	self.AllPorts = os.Args[2:]
	self.IsLeader = false
	self.SetNextNeighbor()
	self.SetPID()
	self.PrintState()
	self.Connections = make(map[string]net.Conn)

	go processLogic(self)

	// keep alive
	fmt.Scanln()
}