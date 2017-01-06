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

func extractCommand (message string) (filename string, tag string, data string) {
	parts := strings.Split(message, ",")
	_, filename = GetKeyValuePair(parts[1])
	_, tag = GetKeyValuePair(parts[2])
	if len(parts) > 3 {
		for i := 3; i < len(parts); i ++ {
			data = data + parts[i]
			if i < len(parts) - 1 {
				data = data + ","
			}
		}
	} else {
		data = ""
	}
	return filename, tag, data
}

func validateCommand (message string) bool {
	parts := strings.Split(message, ",")
	if strings.Index(parts[0], "action=") != -1 {
		return true
	}
	return false
}

func registerHandleInput(self *State, message string) {
	if !validateCommand(message) {
		return
	}
	filename, tag, data := extractCommand(message)
	println("Am extras:   " + filename + "  " + tag + "  " + data)
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(state *State, message string) {
	//err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	//check(err)
}
