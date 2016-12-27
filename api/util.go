package api

import (
	"math/rand"
	"time"
	"runtime"
	"reflect"
	"strings"
)

/*
This function generates a string containing random letters or numbers of a given length
 */
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
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
	equalIndex := strings.Split(message, "=")
	key = equalIndex[0]
	value = equalIndex[1]
	return key, value
}