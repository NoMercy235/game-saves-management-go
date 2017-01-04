package api

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
func registerHandleInput(self *State, message string) {

}
