package api

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"strings"
)

/*
This function is used to write anything to a file. It can write either saves or logs
 */
func Write(self *State, command Command, directory string) {
	message := ""
	path := ""
	if directory == "files" {
		message = command.MakeSave()
		path = directory + "/" + command.Filename
	} else
	if directory == "logs"{
		message = command.ToString()
		path = directory + "/logs"
	}
	CreateFile(path)
	go WriteFile(message, path)
}

/*
This function reads the content from a file and returns the data from the tag that is being searched
 */
func Read(self *State, command Command, directory string) {
	fileData := ReadFile(directory + "/" + command.Filename)
	save := getTagInFileData(command, fileData)
	if save == "" {
		save = "Save not found"
	}
	println("Sent to user: " + save)
}

/*
This function checks if the directory for the files has been created, and it it hasn't been,
it created it with the 0644 (READ and WRITE only for the owner) permission.
Due to the last modifications, the path is actually the filepath, not the directory's, so it needs to go back one level.
 */
func checkDirectory(path string){
	lastSlash := strings.LastIndex(path, "/")
	mainPath := path[:lastSlash]
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		os.Mkdir(mainPath, 0644)
	}
}

/*
Simply checks if the given path is a file and if it exists
 */
func CheckFile(fileName string) bool{
	if _, err := os.Stat(FILES_PATH + fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
This function first checks to see if the directory exists before attempting to do any action
there and then creates a file with the name filename (no extension) only if it does not exist
already
 */
func CreateFile(path string) {
	checkDirectory(FILES_PATH + path)
	// detect if file exists
	var _, err = os.Stat(FILES_PATH + path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(FILES_PATH + path)
		if hasError(err) == true {
			return
		}
		defer file.Close()
	}
}

/*
this function simply writes the message to the designated file
 */
func WriteFile(message string, path string) {

	// open file using 0644 (see above) permission
	if !CheckFile(path) || path == "" {
		return
	}

	// Mutual Exclusion
	MUTEX.Lock()
	
	var file, err = os.OpenFile(FILES_PATH + path, os.O_APPEND|os.O_WRONLY, 0644)
	if hasError(err) == true {
		println("EROARE")
		return
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if hasError(err) == true {
		return
	}
	// save changes
	err = file.Sync()
	if hasError(err) == true {
		return
	}

	MUTEX.Unlock()
}

/*
this function reads the data from a file and returns it. currently, it can't read more
than 1024 characters, but that should do it for our case.
 */
func ReadFile(path string) (string){
	var file, err = os.OpenFile(FILES_PATH + path, os.O_RDWR, 0644)
	if hasError(err) == true {
		return ""
	}
	defer file.Close()

	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if hasError(err) == true {
				return ""
			}
		}
		if n == 0 {
			break
		}
	}
	if hasError(err) == true {
		return ""
	}

	strlen := bytes.IndexByte(text, 0)
	if strlen == -1 || strlen > len(text) {
		print("ERROR!!! read from file " + file.Name())
		return ""
	}
	return string(text[:strlen])
}

/*
This function is basically useless unless someone actually would want to delete their
saves. So I'll leave it here just in case
IMPORTANT !!! if you ever use this, it needs to be updated
 */
func DeleteFile(command Command, directory string) {
	var err = os.Remove(FILES_PATH + directory + "/" + command.Filename)
	if hasError(err) == true {
		return
	}
}

func hasError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}