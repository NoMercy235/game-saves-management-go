package api
import (
	"math/rand"
	"time"
	"strconv"
	"bytes"
)

/*
Think of a standard format of sending saves. We shouldn't bother with JSON. just send a string and parse it
with indexOf or something of the sort.

We will also have two actions: write and read.

Then message should look something like this:
write,filename:'first-player.save',life:100, money:0,tag:'first save'
read,filename:'first-player.save',tag:'first save'

Then, implement the following logic:
 - if the process is not the leader: generate an action at a random interval with random data
 - if it's a leader: this one is more complicated but for now just make it execute the action
 (write to file and send something as ACK back or read and return the result)

 P.S. the function name is just... to be there. you can organize it however you want
 */

var lastCommand string
func GenerateInput(self *State) {
	println("*** Starting to generate input to send to " + self.LeaderPort + " ! ***")
	for {
		time.Sleep(MESSAGE_TIME)
		command := ""
		// If the leader is lost, then store the last action until a new one is found
		// and retry sending it afterwards
		if lastCommand == "" {
			command = generateCommand(self)
		} else {
			command = lastCommand
		}
		if(self.LeaderPort == "") {
			lastCommand = command
			break
		}
		//println("I am about to send: " + GetFriendlyCommand(self, command))
		Send(self, self.LeaderPort, command)
		lastCommand = ""
	}
}

/*
This is a private function that handles generation of command such as:
write,filename='[randomString]',tag='[randomString]',life=100, money=0
read,filename='[randomString]',tag='[randomString]'
 */
func generateCommand(self *State) string{
	if(self.LeaderPort == ""){
		return ""
	}
	actions := [2]string{"action=write", "action=read"}
	var command bytes.Buffer
	index := 0;
	rnd := rand.Float64()
	if rnd < 0.5 {
		index = 0
	} else {
		index = 1
	}
	command.WriteString(actions[index])
	command.WriteString(",filename=")
	command.WriteString(RandomString(20))
	command.WriteString(",tag=")
	command.WriteString(RandomString(20))
	if(rnd < 0.5){
		command.WriteString(",life=")
		command.WriteString(strconv.Itoa(rand.Intn(20)))
		command.WriteString(",money=")
		command.WriteString(strconv.Itoa(rand.Intn(20)))
	}
	return command.String()
}

