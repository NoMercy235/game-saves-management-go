package api

import (
	"strings"
	"time"
	"strconv"
)

func RegisterHandleInputCallbacks(self *State) {
	self.RegisterCallback(registerHandleInput)
}

/*
This function will take care of all the processing of the received messages. Here, the mutual
exclusion algorithm should be used for data integrity and such. It should also (probably, I
don't know if it's possible) use the synchronized clock (the algorithm is still not implemented)
to process requests. But this is for another time.

Basically what it has to do:
- check if the message is a read or write action. If it is something else (not possible, but
best check for it), simply return
- if it's a read, should read the content of the given filename at the given tag and Send() it
to the process that requested it (I wonder how that's done :-?)
- if it's a write, (use the mutual exclusion to lock the file and write to it and...) simply
 write the contents of the message to the file as intended (for now)

 IMPORTANT!!!
 This is also where we should somehow implement the two-phase commit thingy. We've talked about
 this, but if you don't remember, check the courses.
 */

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
	parts := strings.Split(message, ",")
	if len(parts) < 2 {
		return false
	}
	if strings.Index(parts[0], "source=") != -1 && strings.Index(parts[1], "action=") != -1 {
		return true
	}
	return false
}

/*
This function updates the commands queue of the leader whenever it receives a command.
 */
func registerHandleInput(self *State, message string) {
	if self.IsLeader == false || !validateCommand(message) {
		return
	}
	command := extractCommand(message)
	updateQueue(self, command)
}

// Check if other process is having an action by asking the leader
func updateQueue(self *State, command Command) {
	var hasAction bool
	for i := 0; i < len(self.CommandsQueue); i++ {
		queueCommand := self.CommandsQueue[i];
		// If command is already taken by another process
		if queueCommand.Filename == command.Filename {
			// Push the action to the queueCommand
			if len(self.CommandsQueue) < COMMAND_QUEUE_LIMIT {
				self.RegisterCommand(command, true)
				hasAction = true
			}
			break
		}
	}

	if !hasAction {
		self.RegisterCommand(command, false)
	}
	println("Current command queue length: " + strconv.Itoa(len(self.CommandsQueue)))
}

/*
This function enters in an infinite loop that tries to execute any command it finds on a leader
 */
func ExecuteCommands(self *State) {
	for {
		time.Sleep(EXECUTE_COMMAND_DELAY)
		if len(self.CommandsQueue) != 0 {
			command := self.PopCommand()

			if command.Action == "write" {
				write(self, command)
			} else {
				read(self, command)
			}
		}
	}
}

func write(self *State, command Command) {
	CreateFile(command)
	go WriteFile(command)
}


func read(self *State, command Command) {
	fileData := ReadFile(command)
	go Send(self, command.SourcePort, getTagInFileData(command, fileData))
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
