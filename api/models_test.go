package api

import (
	"testing"
	"strings"
	"strconv"
	"time"
)

func TestState(t *testing.T) {
	states := []State {}
	var state State
	state.ListenPort = "8081"
	state.AllPorts = []string{"8081", "8082", "8083"}
	states = append(states, state)

	for _, state := range(states) {
		setNextNeighbor := func(t *testing.T, state State) {
			state.SetNextNeighbor()
			if state.SendPort == "" {
				t.Log("Expected state.SendPort not \"\"")
				t.Fail()
			}
		}
		setNextNeighbor(t, state)

		setPID := func(t *testing.T, state State) {
			state.SetPID()
			if state.PID == 0 {
				t.Log("Expected state.PID not \"\"")
				t.Fail()
			}
		}
		setPID(t, state)

		registerCallback := func(t *testing.T, state State) {
			testFunc := func(state *State, message string) {}
			state.RegisterCallback(testFunc)
			if len(state.Callbacks) == 0 {
				t.Log("Expected len(state.Callbacks) == 1")
				t.Log("Got 0")
				t.Fail()
			}
		}
		registerCallback(t, state)

		generateLeaderToken := func(t *testing.T, state State) {
			token := state.GenerateLeaderToken()
			if strings.Index(token, "token=") == -1 || strings.Index(token, "pid=") == -1 {
				t.Log("Malformed or invalid token")
				t.Log("Got: " + token)
				t.Fail()
			}
		}
		generateLeaderToken(t, state)

		registerCommand  := func(t *testing.T, state State) {
			var command Command
			command.Action = "write"
			state.RegisterCommand(command, true)
			if len(state.CommandsQueue) == 0 {
				t.Log("Expected len(state.CommandsQueue) == 1")
				t.Log("Got 0")
				t.Fail()
			}

			command.Action = "read"
			state.RegisterCommand(command, false)
			if state.CommandsQueue[0].Action != "read" {
				t.Log("Expected first command to have action=read")
				t.Log("Got: " + state.CommandsQueue[0].Action)
				t.Fail()
			}
		}
		registerCommand(t, state)

		popCommand := func(t *testing.T, state State) {
			var command Command
			command.Action = "write"
			state.RegisterCommand(command, true)
			res := state.PopCommand()

			if res.Action != "write" {
				t.Log("Expected received command to have action=write")
				t.Log("Got: " + state.CommandsQueue[0].Action)
				t.Fail()
			}

			if len(state.Callbacks) != 0 {
				t.Log("Expected len(state.Callbacks) == 0")
				t.Log("Got: " + strconv.Itoa(len(state.Callbacks)))
				t.Fail()
			}
		}
		popCommand(t, state)
	}
}

func TestCommand(t *testing.T) {
	commands := []Command {}
	var command Command
	command.SourcePort = "8081"
	command.Action = "write"
	command.Filename = "file1"
	command.Filename = "tag1"
	command.Data.Life = "100"
	command.Data.Money = "10"
	commands = append(commands, command)

	for _, command := range commands {
		gameDataToString := func (t *testing.T, command Command) {
			result := command.Data.ToString()
			if result != "life=100&money=10" {
				t.Log("Expected: life=100&money=10")
				t.Log("Got: " + result)
				t.Fail()
			}
		}
		gameDataToString(t, command)

		commandToString := func (t *testing.T, command Command) {
			result := command.ToString()
			if result != ("source=" + command.SourcePort + ",action=" + command.Action + ",filename=" + command.Filename +
				",tag=" + command.Tag + "," + command.Data.ToString()) {
				t.Log("Expected: " + ("source=" + command.SourcePort + ",action=" + command.Action + ",filename=" + command.Filename +
					",tag=" + command.Tag + "," + command.Data.ToString()))
				t.Log("Got: " + result)
				t.Fail()
			}
		}
		commandToString(t, command)

		makeSave := func (t *testing.T, command Command) {
			result := command.MakeSave()
			if result != (command.Tag + " : " + command.Data.ToString()) {
				t.Log("Expected: " + command.Tag + " : " + command.Data.ToString())
				t.Log("Got: " + result)
				t.Fail()
			}
		}
		makeSave(t, command)
	}
}

func TestInternalClock(t *testing.T) {
	clocks := []InternalClock{}
	var clock InternalClock
	clock.ServerRtt = 100 * time.Millisecond
	clocks = append(clocks, clock)

	for _, clock := range clocks {
		setRealTime := func (t *testing.T, clock InternalClock) {
			clock.SetRealTime()
			if clock.Synchronized == false {
				t.Log("Expected the clock to be synchronized!")
				t.Fail()
			}
		}
		setRealTime(t, clock)
	}
}


func TestProposition(t *testing.T) {
	var self State
	self.ListenPort = "8081"

	var command Command
	command.Action = "write"

	propositions := []Proposition{}
	var proposition Proposition
	//proposition.SetProposalIndex(self)
	propositions = append(propositions, proposition)

	for _, proposition := range propositions {
		isEmpty := func (t *testing.T, proposition Proposition) {
			result := proposition.IsEmpty()
			if !result {
				t.Log("Expected proposition to be empty")
				t.Fail()
			}
		}
		isEmpty(t, proposition)

		setProposalIndex := func (t *testing.T, proposition Proposition) {
			proposition.SetProposalIndex(self)
			if proposition.Index != "1.8081" {
				t.Log("Expected index to be 1.8081")
				t.Log("Got: " + proposition.Index)
				t.Fail()
			}
		}
		setProposalIndex(t, proposition)
	}
}
