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

func registerOnValueProposed(self *State, message string) {
	if !validatePropose(message) {
		return
	}

	index := message[:strings.Index(message, "=")]
	firstBracket := strings.Index(message, "[")
	secondBracket := strings.Index(message, "]")
	command := extractCommand(message[(firstBracket + 1):secondBracket])


	if self.IsLeader {
		if self.Proposition.ProposedValue.IsEmpty() {
			self.Proposition.Index = index
			self.Proposition.ProposedValue.CopyFromCommand(command)
			self.Proposition.Votes = 2
		} else {
			if CompareIndex(self.Proposition.Index, index) == 0 {
				self.Proposition.Votes = self.Proposition.Votes + 1
			} else
			if CompareIndex(self.Proposition.Index, index) == 1 {
				self.Proposition.Index = index
				self.Proposition.ProposedValue.CopyFromCommand(command)
				self.Proposition.Votes = 2
			}
		}

		if self.Proposition.ChosenValue.IsEmpty() && self.Proposition.Votes >= ((len(self.AllPorts) / 2) + 1) {
			self.Proposition.ChosenValue.CopyFromCommand(self.Proposition.ProposedValue)
			println("*** Value " + self.Proposition.Index + " has been chosen! ***")
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
		if CompareIndex(self.Proposition.Index, index) == 1 || self.Proposition.IsEmpty() {
			self.Proposition.Index = index
			self.Proposition.ProposedValue.CopyFromCommand(command)
			go Send(self, self.LeaderPort, self.GetProposal(command), true)
		}
	}

}

func validateChosen(message string) bool {
	match, _ := regexp.MatchString("chosen=", message)
	match2, _ := regexp.MatchString("index=", message)
	return match && match2
}

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