package api

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
	"reflect"
)

func PrintState(self State) {
	fmt.Printf("State config: \nSend port: %s\nListen port: %s \nNetwork config: %s\n\n\n", self.SendPort, self.ListenPort, self.AllPorts)
}

/*
This function gets a state and populates the SendPort property based on the state's place in the topology array
 */
func GetNextNeighbor(self *State) (string) {
	for index, port := range self.AllPorts {
		if self.ListenPort == port {
			neighborIndex := -1
			if index + 1 >= len(self.AllPorts) {
				neighborIndex = 0
			} else {
				neighborIndex = index + 1
			}
			return self.AllPorts[neighborIndex]
		}
	}
	return ""
}
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}