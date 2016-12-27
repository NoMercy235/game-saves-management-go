package api

import (
	"strings"
	"strconv"
)

//import "time"


/*
This file should do the following:

 If i am a leader:
 - respond to others that are trying to see if I'm still alive

 If i am not a leader:
 - check to see if the leader still lives, and if it is not, initiate the leader election algorithm


 Leader election algorithm:
 We will use the token technique, visualizing the topology in a ring. The process that has noticed that the leader
 is gone, will initiate this algorithm. It will instantiate a random token (we should make a struct for this), and will
 fill a property of it with it PID then pass it to the next in circle (state.SendPort).
 The token will get passed like that and every process will fill the PID property with their PID *only* if their PID
 is grater than what was already there.
 When a process gets the token and find the same PID as its own, it then becomes a leader and sends a message to all the
 others to notify them (they all set the state.LeaderPort property

 IMPORTANT: when starting the app, all the processes will have "" in the LeaderPort. Which means that all of them will
 initiate the algorithm, leading to many tokens travelling through the topology. Should use a lot of logs to see how
 exactly all of this is happening
 */

/*
This function simply registers the registerTokenCallback on the current state. Could've done it directly, but, in case
there'll be more callbacks here, it's better to leave a single function to do the work.
 */
func RegisterTokenCallback(self *State){
	self.RegisterCallback(registerTokenCallback)
}

/*
This is the callback that handles the leader election algorithm. It should handle the case where the message received
is not related to the election algorithm (it's not of type token:[randomString],pid:[pid]). We could try some
token validation (instead of randomly generating them, we could (encrypt) the PID of the receiving token and send it,
ot send a text ecrypted using the PID of the receving process as a key. IDK :-? )
 */
func registerTokenCallback(self *State, message string)  {
	messageParts := strings.Split(message, ",")
	pid := messageParts[1][4:]
	intPid, _ := strconv.Atoi(pid)
	if self.PID > intPid {
		pid = strconv.Itoa(self.PID)
	} else
	if self.PID == intPid {
		// this process becomes the leader. here, the algorithm should stop
		// set the LeaderPort and IsLeader fields
		// make sure everyone knows who's the leader
		println("**** I, " + self.ListenPort + ", am the leader! ***")
	}
	newToken := messageParts[0] + ",pid=" + pid
	go Send(self, newToken)
}

/*
This function is called from the processLogic function of the app.go file and should handle the logic of starting the
leader election algorithm whenever there is no leader present. That includes the following cases:
- The application just started and there's no leader present
- The leader hasn't responded for some time, is considered down, and a new one should be elected
 */
func CheckLeader(self *State){
	/*
	This is a hack and MUST be changed. It means that, when the application starts, it will initiate the leader
	election algorithm only if the process is listening on port 8081. Could work without that, but it ust be tested
	 */
	if self.LeaderPort == "" && self.ListenPort == "8081" {
		println("Starting leader algorithm from: " + self.ListenPort)
		chooseLeader(self)
	}
}

/*
This function simply sends a message. If there won't be any other logic necessary, it could very well reside in the
CheckLeader function
 */
func chooseLeader(self *State){
	//Send(self, self.GenerateLeaderToken())
	result := Send(self, self.GenerateLeaderToken())
	if result == -1 {
		println("\n *** Error occured when sending from [" + self.ListenPort + "] to [" + self.SendPort + "]! ***")
	}
}
