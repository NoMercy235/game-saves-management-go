package api

import (
	"math/rand"
	"runtime"
	"reflect"
	"strings"
	"time"
)

/*
This function generates a string containing random letters or numbers of a given length
 */
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

/*
This function really is useless now, but was helpful during the time Callbacks were implemented so I'm leaving it
here as a tribute.
 */
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetKeyValuePair(message string) (key string, value string) {
	parts := strings.Split(message, "=")
	if len(parts) < 2 {
		return "", ""
	}
	key = parts[0]
	value = parts[1]
	return key, value
}

/*
The original version of a command is too cluttered to display in terminal so use this functions to
get a friendlier version
 */
func GetFriendlyCommand(self *State, command string) (friendlyCommand string) {
	parts := strings.Split(command, ",")
	_, tag := GetKeyValuePair(parts[1])
	friendlyCommand = parts[0] + " on " + tag
	return friendlyCommand
}

/*
This function gets the trailing milliseconds of a date
 */
func GetTrailingMilliseconds(time string) (nr string) {
	index := strings.Index(time, ".")
	if index == -1  {
		nr = ""
		return nr
	}

	time = time[index + 1 :]
	index = strings.Index(time, " ")

	aux := time[:index]
	b := make([]rune, len(aux))
	for i := range b {
		b[i] = '0'
	}
	nr = string(b)
	return nr
}


/*
Change the value of a (usually global) time.
 */
func ChangeTime(target *time.Duration, value int32, timeType time.Duration) {
	*target = time.Duration(rand.Int31n(value)) * timeType + MIN_TIME
}