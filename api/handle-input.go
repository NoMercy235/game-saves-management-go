package api

import (
	"strings"
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
 This is also where we should somehow implement the two-face commit thingy. We've talked about
 this, but if you don't remember, check the courses.
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

func registerHandleInput(self *State, message string) {
	if !validateCommand(message) {
		return
	}
	command := extractCommand(message)
	if command.Action == "write" {
		write(self, command)
	} else {
		read(self, command);
	}
}

func write(state *State, command Command) {
	CreateFile(command)
	WriteFile(command)
}

/*
This won't be used until the commands will not be human readable.
Why? Because a read command will always attempt to read from a random file that does not exist.
Solution? Instead of making random filenames and tags, make an array with some hardcoded
filenames and tags and just pick a random index to choose from that array.
 */
func read(state *State, command Command) {
	// This should be save = .... , but...
	//x := ReadFile(command)
	//print("SHOULD SEND BACK: " + x)
	// after reading the contents of a file, the server should send back the information
	// to the process (port) that requested it.
	// I have yet to discover how to get the port of someone that connected to the
	// server. Need to solve the problem
}
