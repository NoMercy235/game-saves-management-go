package api

import (
	"regexp"
	"strings"
	"strconv"
)

func RegisterPaxosCallbacks(self *State) {
	self.RegisterCallback(registerOnValueProposed)
	self.RegisterCallback(registerOnChosenReceived)
}

/*
This function checks if the current process doesn't have a proposition yet and then proposes a value to all the other
processes.
 */
func ProposeValue(self *State, command Command) {
	if !self.Proposition.ProposedValue.IsEmpty() {
		return
	}
	message := self.GenerateProposal(command)
	println("*** Proposed: " + message + " ***");
	for i := 0; i < len(self.AllPorts); i++ {
		if self.ListenPort != self.AllPorts[i] {
			go Send(self, self.AllPorts[i], message, false)
		}
	}
}

func validatePropose(message string) bool {
	match, _ := regexp.MatchString("([0-9]+).([0-9+])=" + regexp.QuoteMeta("[") + "([a-z]+)", message)
	return match
}

/*
This callback is called whenever a process receives a proposition.
If it's a leader process, it checks to see if the value has already been proposed, and if that's the case, it counts the
number of votes. If the votes exceed 50% + 1 of the current number of processes, the proposed value becomes chosen and
it's sent for synchronization
 */
func registerOnValueProposed(self *State, message string) {
	if !validatePropose(message) {
		return
	}

	// extract the proposed value (command) and the index
	index := message[:strings.Index(message, "=")]
	firstBracket := strings.Index(message, "[")
	secondBracket := strings.Index(message, "]")
	command := extractCommand(message[(firstBracket + 1):secondBracket])


	if self.IsLeader {
		// if the leader never had a proposed value until now, it takes the one it received and sets 2 votes for
		// it (the leader itself and the process that proposed it)
		if self.Proposition.ProposedValue.IsEmpty() {
			self.Proposition.New(index, command)
			self.Proposition.Votes = 2
		} else {
			//if there was already a value proposed, it checks if the proposed value that was sent is equal
			// to the one it already had.
			// if it is, the number of votes is incremented
			// if it's not, and the received index is higher than the one it had, the value updates
			if CompareIndex(self.Proposition.Index, index) == 0 {
				self.Proposition.Votes = self.Proposition.Votes + 1
			} else
			if CompareIndex(self.Proposition.Index, index) == 1 {
				self.Proposition.New(index, command)
				self.Proposition.Votes = 2
			}
		}

		// in the end, if there is no chosen value and the number of votes of the current proposition exceeds
		// 50% + 1 of the current number of processes, the proposition is chosen and sent for synchronization
		if self.Proposition.ChosenValue.IsEmpty() && self.Proposition.Votes >= ((len(self.AllPorts) / 2) + 1) {
			self.Proposition.ChosenValue.CopyFromCommand(self.Proposition.ProposedValue)
			println("*** Value " + self.Proposition.Index + " has been chosen! ***")
			// let all the other processes know
			for i := 0; i < len(self.AllPorts); i ++ {
				if self.AllPorts[i] != self.ListenPort {
					go Send(self, self.AllPorts[i],
						"chosen=[" + self.Proposition.ChosenValue.ToString() + "],index=" + self.Proposition.Index,
						false)
				}
			}
			synchronize(self, index)
		}
	} else
	if !self.IsLeader {
		// if the process is not a leader and the value it received had a higher index, or it didn't have any
		// proposition to begin with, it updates it's own proposition and sends it to the leader
		if CompareIndex(self.Proposition.Index, index) == 1 || self.Proposition.IsEmpty() {
			self.Proposition.New(index, command)
			go Send(self, self.LeaderPort, self.GetProposal(command), true)
		}
	}

}

/*
A simple check function that verifies if a message is related to the chosen value callback
 */
func validateChosen(message string) bool {
	match, _ := regexp.MatchString("chosen=", message)
	match2, _ := regexp.MatchString("index=", message)
	return match && match2
}

/*
This callback will respond from the followers of the leader when they receive the value chosen by the leader.
They will extract it and attempt to synchronize.
 */
func registerOnChosenReceived(self *State, message string) {
	if !validateChosen(message) || self.IsLeader {
		return
	}

	firstBracket := strings.Index(message, "[")
	secondBracket := strings.Index(message, "]")
	self.Proposition.ChosenValue.CopyFromCommand(extractCommand(message[(firstBracket + 1):secondBracket]))

	restString := message[(secondBracket + 1):]
	_, index := GetKeyValuePair(restString)

	synchronize(self, index)
}


/*
This function will write the log and save to the respective files in order to synchronize the process with all the others.
After it's completion, the Proposal object will reset and increment the index value, in order to never have a smaller index
on future proposals.
 */
func synchronize(self *State, index string) {
	println()
	println("*** Synchronized with index: " + index + " ***")
	println("*** And value: " + self.Proposition.ChosenValue.MakeSave() + " ***")
	println()
	Write(self, self.Proposition.ChosenValue, "files")
	Write(self, self.Proposition.ChosenValue, "logs")
	indexParts := strings.Split(index, ".")
	intIndex, _ := strconv.Atoi(indexParts[0])
	self.Proposition.Clear(intIndex)
}