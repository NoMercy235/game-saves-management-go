package api

import (
	"os"
	"strconv"
	"fmt"
	"net"
	"time"
)

// This is the base structure for our application.
// ListenPort - the port on which the current process listens to
// SendPort - the port on which the current process will send to (used in algorithms, like the leader election)
// AllPorts - an array containing the ports of all processes.
// LeaderPort & IsLeader - the port of the leader and a boolean that shows if the process is a leader or not
// Callbacks - an array of functions with a (State, string) signature. Those are used to respond to receive message events
// CommandsQueue - a queue for the commands that will wait the other processes to finish their common action
// PID - the PID of the process
// Connections - a map of all the connections of a process.
// Clock - the clock of the state
// Proposition - used in Paxos to propose and persist values
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
	Clock InternalClock
	Proposition Proposition
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

/*
This function generates a proposal. A new proposal will also increase the index value
The proposal will look like this:
1.8081=[command]
 */
func (this *State) GenerateProposal(command Command) string {
	this.Proposition.SetProposalIndex(this)
	this.Proposition.ProposedValue.CopyFromCommand(command)
	return this.GetProposal(command)
}

/*
Same as GenerateProposal, but it no longer increments the index
 */
func (this *State) GetProposal(command Command) string {
	if this.Proposition.ProposedValue.IsEmpty() {
		this.Proposition.ProposedValue.CopyFromCommand(command)
	}
	return  this.Proposition.Index + "=[" + this.Proposition.ProposedValue.ToString() + "]"
}

/*
Used to register a command in the CommandsQueue.
 */
func (this *State) RegisterCommand(command Command, appendCommand bool) {
	if appendCommand {
		this.CommandsQueue = append(this.CommandsQueue, command)
	} else {
		this.CommandsQueue = append([]Command{command}, this.CommandsQueue...)
	}
}

/*
It gets the first command, removes it from the queue and returns it
IMPORTANT!!! look to the usage of this function. A not persisted popped command is lost forever
 */
func (this *State) PopCommand() (command Command) {
	if len(this.CommandsQueue) < 1 {
		return command
	}
	command = this.CommandsQueue[0]
	this.CommandsQueue = append(this.CommandsQueue[1:])
	return command
}

func (this *State) GenerateDateMessage() string {
	return "source=" + this.ListenPort + ",date=" + time.Now().String()
}

/*********************************    Command Struct      ****************************************/
type GameData struct {
	Life string
	Money string
}

func (this *GameData) ToString() string {
	return "life=" + this.Life + "&money=" + this.Money
}

func (this *GameData) CopyFromGameDate(gameData GameData) {
	this.Life = gameData.Life
	this.Money = gameData.Money
}

func (this *GameData) IsEmpty() bool {
	return this.Life == "" && this.Money == ""
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

func (this *Command) IsEmpty() bool {
	return this.Action == "" && this.Filename == "" && this.SourcePort == "" && this.Tag == "" &&
		&this.Data != nil && this.Data.IsEmpty()
}

func (this *Command) CopyFromCommand(command Command) {
	this.Action = command.Action
	this.Filename = command.Filename
	this.SourcePort = command.SourcePort
	this.Tag = command.Tag
	this.Data.CopyFromGameDate(command.Data)
}

/*********************************    InternalClock Struct      ****************************************/
type InternalClock struct {
	Clock time.Time
	StartedSyncTime time.Time
	ServerRtt time.Duration
	Synchronized bool
}

/*
This function will set the real time of the state based on the RTT it received from the server
and set the Synchronized value to true
 */
func (this *InternalClock) SetRealTime() {
	now := time.Now()
	this.Clock = now.Add(this.ServerRtt)
	if this.Synchronized == false {
		this.Synchronized = true;
	}
	println("----------------------------------------------")
	println("My time is: " + this.Clock.String())
	println("Rtt: " + this.ServerRtt.String())
	println("----------------------------------------------")
}



/*********************************    Proposition      ****************************************/
type Proposition struct {
	Index string
	index int
	ChosenValue Command
	ProposedValue Command
	Votes int
}

func (this *Proposition) IsEmpty() bool {
	return this.Index == "" && this.ProposedValue.IsEmpty()
}

/*
This function generate the proposal index of form: 1.8081
The ListenPort is used to make sure that the indexes will be unique and don't affect the indexes
comparison, unless there is no other way.
 */
func (this *Proposition) SetProposalIndex(self *State) {
	this.index += 1
	this.Index = strconv.Itoa(this.index) + "." + self.ListenPort
}

/*
This function clears the proposition
 */
func (this *Proposition) Clear(index int) {
	this.index = index
	this.ChosenValue = *new(Command)
	this.ProposedValue = *new(Command)
	this.Votes = 0
}

/*
This function updates the proposition with a new one
 */
func (this *Proposition) New(index string, value Command) {
	this.Index = index
	this.ProposedValue.CopyFromCommand(value)
}