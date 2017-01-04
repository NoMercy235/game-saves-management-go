package api

import (
	"strings"
	"strconv"
	"time"
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
func RegisterLeaderCallbacks(self *State){
	self.RegisterCallback(leaderTokenCallback)
	//self.RegisterCallback(hasLeaderCallback)
}

/*
This is the callback that handles the leader election algorithm. It should handle the case where the message received
is not related to the election algorithm (it's not of type token:[randomString],pid:[pid]). We could try some
token validation (instead of randomly generating them, we could (encrypt) the PID of the receiving token and send it,
ot send a text ecrypted using the PID of the receving process as a key. IDK :-? )
 */
func leaderTokenCallback(self *State, message string)  {
	if strings.Index(message, "token") == -1 && strings.Index(message, "pid") == -1 {
		// this message has nothing to do with the leader election algorithm, so we just ignore it
		return;
	} else {
		key, token := GetKeyValuePair(strings.Split(message, ",")[0])
		// The length of the token is hardcoded in the GenerateLeaderToken() function. Should probably make it a
		// constant somewhere
		if key != "token" && len(token) != 10 {
			//message is malformed or invalid
			println("Malformed or invalid message: [" + message + "] on process: " + self.ListenPort)
			return ;
		}
	}

	messageParts := strings.Split(message, ",")
	_, pid := GetKeyValuePair(messageParts[1])
	intPid, _ := strconv.Atoi(pid)
	leaderMsg := ""

	// A leader has been elected, because there is a third part: leader=[leaderPort]
	if len(messageParts) >= 3 {
		key, value := GetKeyValuePair(messageParts[2])
		if key != "" && value != "" && key == "leader" {
			// If I'm the process that become leader
			if self.ListenPort == value {
				println("*** Everyone aknowledged me as the leader! ***")
				// Now, when the HandleInput callback will be called from the
				// Listen() function, I, as the leader will be allowed to take
				// care of whatever came
				return
			}
			if self.LeaderPort != "" {
				// A leader has already been elected
				return
			}
			// If I'm a process that become client, simply acknowledge the leader and send the message forward
			self.LeaderPort = value
			println("*** I, " + self.ListenPort + ", aknowledge " + self.LeaderPort + " as the leader! ***")
			go Send(self, self.SendPort, message)
			// A client will periodically ping the leader to see if it's still alive
			go pingLeader(self)
			// A client will randomly generate input to be sent to the leader for processing
			// In a real application, this operation could be very demanding, and thus, it's better
			// to be taken care of in a client process
			go GenerateInput(self)
		}
		return
	} else
	if self.PID > intPid {
		pid = strconv.Itoa(self.PID)
	} else
	if self.PID == intPid {
		if self.LeaderPort != "" {
			// A leader has already been elected. This is needed when the algorithm
			// is initiated by more processes. In that case, a leader elects itself
			// as the leader more than once.
			return
		}
		// this process becomes the leader. here, the algorithm should stop
		// set the LeaderPort and IsLeader fields
		// make sure everyone knows who's the leader
		println("*** I, " + self.ListenPort + ", am the leader! ***")
		self.IsLeader = true
		self.LeaderPort = self.ListenPort
		leaderMsg = ",leader=" + self.ListenPort
	}
	newToken := messageParts[0] + ",pid=" + pid + leaderMsg
	go Send(self, self.SendPort, newToken)
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
		chooseLeader(self)
	}
}

/*
This function simply sends a message. If there won't be any other logic necessary, it could very well reside in the
CheckLeader function
 */
func chooseLeader(self *State){
	println("Starting leader algorithm from: " + self.ListenPort)
	result := Send(self, self.SendPort, self.GenerateLeaderToken())
	if result == -1 {
		println("\n *** Error occured when sending from [" + self.ListenPort + "] to [" + self.SendPort + "]! ***")
	}
}

/*
This function is called after a leader has been elected and it's role is to check the connectivity
with the leader from time to time and, in case the leader does not respond, it initiates an algorithm
of leader election
 */
func pingLeader(self *State){
	for {
		time.Sleep(PING_TIME)
		result := Send(self, self.LeaderPort, "--- Ping from " + self.ListenPort + " ---")
		if result == -1 {
			println("*** Leader is down ***")
			removeLeader(self)
			chooseLeader(self)
			break
		}
	}
}

/*
This function is called when a leader is declared dead to remove it from the topology array (self.AllPorts),
reset the self.LeaderPort field and choose a new neighbor.
 */
func removeLeader(self *State) {
	self.RemovePort(self.LeaderPort)
	self.LeaderPort = ""
	self.SetNextNeighbor()
	self.PrintState()
}