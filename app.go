package main

import (
	"os"
	"fmt"
	"./api"
)


func processLogic(self api.State) {
	go api.CheckLeader(self)
	go api.GenerateInput(self)
	//call the other functions as they are made
}

func main() {
	if len(os.Args) < 3 {
		println("Wrong usage")
		return
	}

	var self api.State;
	self.ListenPort = os.Args[1]
	self.AllPorts = os.Args[2:]
	self.IsLeader = false
	self.SendPort = api.GetNextNeighbor(self)
	api.PrintState(self)

	go processLogic(self)

	// keep alive
	fmt.Scanln()
}