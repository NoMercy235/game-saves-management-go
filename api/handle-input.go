package api

import (
	"strings"
	"time"
	"strconv"
	//"regexp"
)

func RegisterHandleInputCallbacks(self *State) {
	self.RegisterCallback(registerHandleInput)
}

/*
This function parses a command from a given string
 */
func extractCommand (message string) (command Command) {
	parts := strings.Split(message, ",")
	_, command.SourcePort = GetKeyValuePair(parts[0])
	_, command.Action = GetKeyValuePair(parts[1])
	_, command.Filename = GetKeyValuePair(parts[2])
	_, command.Tag = GetKeyValuePair(parts[3])
	if len(parts) > 4 {
		extras := strings.Split(parts[4], "&")
		_, command.Data.Life = GetKeyValuePair(extras[0])
		_, command.Data.Money = GetKeyValuePair(extras[1])
	}
	return command
}

/*
This function checks to see if a given string matches the pattern of a command
 */
func validateCommand (message string) bool {
	// TODO maybe use regexp
	parts := strings.Split(message, ",")
	if len(parts) < 2 || strings.Index(message, "[") != -1 || strings.Index(message, "]") != -1 {
		return false
	}
	if strings.Index(parts[0], "source=") != -1 && strings.Index(parts[1], "action=") != -1 {
		return true
	}
	return false
}

/*
This function updates the commands queue of a process whenever it receives a command.
 */
func registerHandleInput(self *State, message string) {
	if !validateCommand(message) {
		return
	}
	command := extractCommand(message)
	UpdateQueue(self, command)
}

/*
This function updates the command's queue with a received command. If the command wants to make
an action
If the received command is a read, it first checks to see if there's any other pending command
on the same file that wants to write on it, and if there is, it adds the read command at the
end of the queue. If it's not, it pushes it upfront
A write command will always be pushed to the front.

IMPORTANT!!! This is probably not a good idea since a write command can be starved
A new write command should be inserted after the last write in the queue, regardless of the
existence of any read commands
 */
func UpdateQueue(self *State, command Command) {
	var hasAction bool
	if command.Action == "read" {
		for i := 0; i < len(self.CommandsQueue); i++ {
			queueCommand := self.CommandsQueue[i];
			if queueCommand.Filename == command.Filename && queueCommand.Action == "write" {
				if len(self.CommandsQueue) < COMMAND_QUEUE_LIMIT {
					self.RegisterCommand(command, true)
					hasAction = true
				}
				break
			}
		}
	}

	if !hasAction {
		self.RegisterCommand(command, false)
	}
	println("Current command queue length: " + strconv.Itoa(len(self.CommandsQueue)))
}

/*
This function enters in an infinite loop that tries to execute any command it finds on a leader
if a read command is received, it executes locally, since the system should be synchronized
if a write command is received, it uses the Paxos protocol to try and make it a chosen value
and persist it

IMPORTANT!!! if a command is received, but Paxos decides that another one should be persisted,
that command is lost forever. Should fix this, probably.
 */
func ExecuteCommands(self *State) {
	for {
		time.Sleep(EXECUTE_COMMAND_DELAY)
		if len(self.CommandsQueue) != 0 {
			command := self.PopCommand()

			if command.Action == "write"  {
				go ProposeValue(self, command)
			} else {
				go Read(self, command, "files")
			}
		}
	}
}

/*
This function returns the line of a file that corresponds to a given tag
 */
func getTagInFileData(command Command, fileData string) (string) {
	tagIndex := strings.Index(fileData, command.Tag)
	if tagIndex == -1 {
		return ""
	}
	restOfFile := fileData[tagIndex:]
	endIndex := strings.Index(restOfFile, "\n")
	return restOfFile[:endIndex]
}