package api

import (
	"time"
	"strings"
	"regexp"
)

/*
This file should implement the logic to synchronize the clocks using one of the methods we used for the lab problems
There should be another field in the State structure called 'Clock' which will be the variable to be synchronized.

Logic:
- When a leader is present (lock the code with a while(self.LeaderPort == "") { do nothing; } . but this might be
a bad idea) use one of the known algorithms to synchronize the clock with the one on the server.

P.S. Might be hard to use dates for the clock, so maybe stick with integers? (must find a way to increment them,
maybe make a 'job' and launch it with a go routine to increment every second)
 */

func RegisterClockSyncCallbacks(self *State){
	//self.RegisterCallback(clockSyncCallback)
	self.RegisterCallback(clockReceivedSyncCallback)
	self.RegisterCallback(clockLeaderSyncCallback)
}

/*
This function will attempt to sync the clock with the leader every CLOCK_SYNC_TIME seconds
 */
func StartClockSync(self *State) {
	for {
		time.Sleep(CLOCK_SYNC_TIME)
		if self.LeaderPort == "" || self.IsLeader == true {
			continue
		}
		go Send(self, self.LeaderPort, generateDateMessage(self), false)
	}
}

func validateDate(message string) bool {
	match, _ := regexp.MatchString("([0-9]+)-([0-9+])-([0-9+])", message)
	return match
}

/*
This function is called whenever a process received the clock from the server. It checks the RTT and updates the
received clock accordingly
 */
func clockReceivedSyncCallback(self *State, message string) {
	if self.LeaderPort == "" || self.IsLeader == true || !validateDate(message) {
		return
	}
	layout :=  "2006-01-02 15:04:05." + GetTrailingMilliseconds(message) + " -0700 MST"
	clock, err := time.Parse(layout, message)
	if err != nil {
		self.Clock.Synchronized = false
		return
	}
	self.Clock.ServerRtt = time.Since(clock)
	self.Clock.SetRealTime()
}

/*
This function is called whenever the leader receives a clock sync request from the other processes.
It sends back it's own time
 */
func clockLeaderSyncCallback(self *State, message string)  {
	if self.IsLeader == false || !validateDate(message) {
		return
	}
	parts := strings.Split(message, ",")
	if len(parts) != 2 || strings.Index(parts[0], "source=") == -1 || strings.Index(parts[1], "date=") == -1 {
		// message not related to clock sync
		return
	}
	_, sourcePort := GetKeyValuePair(parts[0])
	go Send(self, sourcePort, time.Now().String(), false)
}

func generateDateMessage(self *State) (message string) {
	return "source=" + self.ListenPort + ",date=" + time.Now().String()
}