package api

import (
	"testing"
	"os"
	"strings"
)

/*
This is a small integration test that will go through the entire functionality offered by
file-api.go
 */
func TestFileApi(t *testing.T){
	cases := []Command {
		{ "", "", "testFile", "myTag", GameData{ "10", "100" } },
	}

	var cwd, _ = os.Getwd()
	lastSlash := strings.LastIndex(cwd, "\\")
	// the location will be the a new directory on the same level as the usual tmp
	FILES_PATH = cwd[:lastSlash] + "/test_tmp/"

	for _, c := range cases {
		createFile := func(t *testing.T) {
			CreateFile(c.Filename)
			result := CheckFile(c.Filename)
			if result == false {
				t.Log("Expected " + c.Filename + " to exist!")
				t.Fail()
			}
		}
		createFile(t)

		writeReadFile := func(t *testing.T) {
			WriteFile(c.MakeSave(), c.Filename)
			result := ReadFile(c.Filename)
			if result == "" {
				t.Log("Expected " + c.MakeSave())
				t.Fail()
			}
			if result[:(len(result) - 1)] != c.MakeSave() {
				t.Log("Write failed.")
				t.Log("Expected: " + c.MakeSave())
				t.Log("Got: " + result)
				t.Fail()
			}
		}
		writeReadFile(t)

		deleteFile := func(t *testing.T) {
			DeleteFile(c.Filename)
		}
		deleteFile(t)

		os.Remove(FILES_PATH) // delete the directory
	}
}
