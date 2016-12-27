package api

import (
	"os"
	"strconv"
	"fmt"
	"net"
)

type State struct {
	ListenPort string
	SendPort string
	AllPorts []string
	LeaderPort string
	IsLeader bool
	Callbacks []func(self *State, message string)
	PID int
	TokenChan chan string
	SendConn net.Conn
}


func (this *State) PrintState() {
	fmt.Printf("State config: \nPID: %d \nSend port: %s\nListen port: %s \nNetwork config: %s\n\n\n", this.PID, this.SendPort, this.ListenPort, this.AllPorts)
}

/*
This function gets a state and populates the SendPort property based on the state's place in the topology array
 */
func (this *State) SetNextNeighbor() {
	for index, port := range this.AllPorts {
		if this.ListenPort == port {
			neighborIndex := -1
			if index + 1 >= len(this.AllPorts) {
				neighborIndex = 0
			} else {
				neighborIndex = index + 1
			}
			this.SendPort = this.AllPorts[neighborIndex]
			return
		}
	}
	return
}

func (this *State) SetPID() {
	this.PID = os.Getpid()
}

func (this *State) RegisterCallback(function func(self *State, message string)) {
	this.Callbacks = append(this.Callbacks, function)
}

func (this *State) GenerateLeaderToken() (string) {
	token := "token:" + RandomString(10) + ",pid:" + strconv.Itoa(this.PID)
	return token
}