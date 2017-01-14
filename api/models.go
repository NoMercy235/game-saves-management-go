package api

import (
	"os"
	"strconv"
	"fmt"
	"net"
)

// This is the base structure for our application.
// ListenPort - the port on which the current process listens to
// SendPort - the port on which the current process will send to (used in algorithms, like the leader election)
// AllPorts - an array containing the ports of all processes.
// LeaderPort & IsLeader - the port of the leader and a boolean that shows if the process is a leader or not
// Callbacks - an array of functions with a (State, string) signature. Those are used to respond to receive message events
// CommandsQueue - a queue for the commands that will wait the other processes to finish their common action
// PID - the PID of the process
// SendConn - the connection used to send messages. Needed to be cached. maybe find a better (more sustainable) solution
type State struct {
	ListenPort string
	SendPort string
	AllPorts []string
	LeaderPort string
	IsLeader bool
	Callbacks []func(self *State, message string)
	CommandsQueue []Command
	PID int
	Connections map[string]net.Conn
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

/*
This function should be used whenever you want to register callbacks that will respond to a message received event. When
that event happens, all callbacks will be called with the state and the message received
 */
func (this *State) RegisterCallback(function func(self *State, message string)) {
	this.Callbacks = append(this.Callbacks, function)
}


/*
This function is used in the leader election algorithm to generate a message containing  a random token (don't even
know if that's needed) along with the PID of the current process
 */
func (this *State) GenerateLeaderToken() (string) {
	token := "token=" + RandomString(10) + ",pid=" + strconv.Itoa(this.PID)
	return token
}


/*
This function removes the port given as a parameter from the array of AllPorts of a state
 */
func (this *State) RemovePort(port string) {
	for i := 0; i < len(this.AllPorts); i++ {
		if (port == this.AllPorts[i]) {
			this.AllPorts = append(this.AllPorts[:i], this.AllPorts[i + 1:]...)
			break;
		}
	}
}


func (this *State) PopCommand() (command Command) {
	if len(this.CommandsQueue) < 1 {
		return command
	}
	command = this.CommandsQueue[0]
	this.CommandsQueue = append(this.CommandsQueue[1:])
	return command
}


/*********************************    Command Struct      ****************************************/
type GameData struct {
	Life string
	Money string
}

func (this *GameData) ToString() string {
	return "life=" + this.Life + "&money=" + this.Money
}

type Command struct {
	SourcePort string
	Action string
	Filename string
	Tag string
	Data GameData
}

/*
This function generates a save from a command
 */
func (this *Command) MakeSave() string {
	return this.Tag + " : " + this.Data.ToString()
}

func (this *Command) ToString() string {
	firstPart := "source=" + this.SourcePort + ",action=" + this.Action + ",filename=" + this.Filename +
		",tag=" + this.Tag
	if this.Data.Life != "" && this.Data.Money != "" {
		return firstPart + "," + this.Data.ToString()
	}
	return firstPart
}