package api

import (
	"fmt"
	"os"
	"io"
	"bytes"
)

/*
This function checks if the directory for the files has been created, and it it hasn't been,
it created it with the 0644 (READ and WRITE only for the owner) permission.
 */
func checkDirectory(){
	if _, err := os.Stat(FILES_PATH); os.IsNotExist(err) {
		os.Mkdir(FILES_PATH, 0644)
	}
}

/*
This function first checks to see if the /tmp directory exists before attempting to do any action
there and then creates a file with the name filename (no extension) only if it does not exist
already
 */
func CreateFile(command Command) {
	checkDirectory()
	// detect if file exists
	var _, err = os.Stat(FILES_PATH + command.Filename)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(FILES_PATH + command.Filename)
		checkError(err)
		defer file.Close()
	}
}

/*
this function simply writes the save to the designated file
 */
func WriteFile(command Command) {
	// open file using 0644 (see above) permission
	var file, err = os.OpenFile(FILES_PATH + command.Filename, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	_, err = file.WriteString(command.MakeSave() + "\n")
	checkError(err)

	// save changes
	err = file.Sync()
	checkError(err)
}

/*
this function reads the data from a file and returns it. currently, it can't read more
than 1024 characters, but that should do it for our case.
 */
func ReadFile(command Command) (string){
	var file, err = os.OpenFile(FILES_PATH + command.Filename, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			checkError(err)
		}
		if n == 0 {
			break
		}
	}
	checkError(err)

	strlen := bytes.IndexByte(text, 0)
	return string(text[:strlen])
}

/*
This function is basically useless unless someone actually would want to delete their
saves. So I'll leave it here just in case
 */
func DeleteFile(command Command) {
	var err = os.Remove(FILES_PATH + command.Filename)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}