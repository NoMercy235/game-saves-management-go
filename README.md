To run this, go to cmd/bash/any editor and run:

./start.bat

To have more processes, edit the start.bat file and add the required lines;
Ex:

start cmd /c "go run app.go" 8081 8081 8082 8083
start cmd /c "go run app.go" 8082 8081 8082 8083
start cmd /c "go run app.go" 8083 8081 8082 8083
pause


This starts 3 processes. The first argument is the listening port and the rest will be placed inside an array representing the processes of the app. They must define all the ports that the other process will use to listen to.


If a process hangs and the ports become unusable, use see.bat and kill.bat to remove them.
inside the kill.bat file, add as many lines as necessary to kill every process on every port.
Ex:

FOR /F "tokens=5 delims= " %%P IN ('netstat -a -n -o ^| findstr :8083') DO TaskKill.exe /PID %%P /F

This kills every process using port 8083.


For testing:
- cd to the api folder
- go test to run every file that ends with *_test.go and execute every function inside that
 starts with TestFunctionName (t *testing.T) { ... }
 // or TestFoo, TestAnything. It just has to start with Test
- go test -cover - for seeing th percentage of coverage
// the following does not work, but the docs said it should :-?
- if you want to see which parts are not covered use this:
go test -coverprofile=coverage.out // or any other name
go tool cover -html=coverage.out
