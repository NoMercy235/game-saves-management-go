To run this, go to cmd/bash/any editor and run:

./start.bat

To have more processes, edit the start.bat file and add the required lines;
Ex:

start cmd /c "go run app.go" 8081 8082
start cmd /c "go run app.go" 8082 8083
start cmd /c "go run app.go" 8083 8081
pause


This starts 3 processes. The first argument is the listening port and the second is the sending port. They must make a complete circle


If a process hangs and the ports become unusable, use see.bat and kill.bat to remove them.
inside the kill.bat file, add as many lines as necessary to kill every process on every port.
Ex:

FOR /F "tokens=5 delims= " %%P IN ('netstat -a -n -o ^| findstr :8083') DO TaskKill.exe /PID %%P /F

This kills every process using port 8083.

