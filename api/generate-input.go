package api

func GenerateInput(self State) {
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
}
