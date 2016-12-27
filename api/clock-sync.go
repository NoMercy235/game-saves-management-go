package api

/*
This file should implement the logic to synchronize the clocks using one of the methods we used for the lab problems
There should be another field in the State structure called 'Clock' which will be the variable to be synchronized.

Logic:
- When a leader is present (lock the code with a while(self.LeaderPort == "") { do nothing; } . but this might be
a bad idea) use one of the known algorithms to synchronize the clock with the one on the server.

P.S. Might be hard to use dates for the clock, so maybe stick with integers? (must find a way to increment them,
maybe make a 'job' and launch it with a go routine to increment every second) 
 */
