package api

import (
	"testing"
	"time"
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