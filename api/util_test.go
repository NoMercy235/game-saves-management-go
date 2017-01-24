package api

import (
	"testing"
	"time"
	"strconv"
)

/*
It should test weather or not the GetKeyValuePair function returns an expected result
and doesn't crash on bad inputs.
 */
func TestGetKeyValuePair(t *testing.T) {
	cases := []struct {data, key, value string} {
		{ "myKey=myValue", "myKey", "myValue" },
		{ "process=me", "process", "me" },
		{ "", "", "" },
		{ "anotherfail", "", "" },
	}

	for _, c := range cases {
		key, val := GetKeyValuePair(c.data)
		if key != c.key || val != c.value {
			t.Log("Expected: key=" + c.key + " and value=" + c.value)
			t.Log("Got: key=" + key + " and val=" + val)
			t.Fail()
		}
	}
}


func TestGetTrailingMilliseconds(t *testing.T) {
	result := GetTrailingMilliseconds(time.Now().String())
	if result == "" {
		t.Log("Expected a number. Got an empty string!")
		t.Fail()
	}
}

func TestCompareIndex(t *testing.T) {
	type Case struct {
		firstIndex string
		secondIndex string
		expected int
	}
	cases := []Case {
		{ "1.8081", "1.8082", 1 },
		{ "3.8081", "1.8082", -1 },
		{ "1.8082", "1.8082", 0 },
		{ "1.808.2", "1.8082", -2 },
	}
	for _, c := range cases {
		result := CompareIndex(c.firstIndex, c.secondIndex)
		if result != c.expected {
			t.Log("Expected: " + strconv.Itoa(c.expected))
			t.Log("Got: " + strconv.Itoa(result))
		}
	}
}