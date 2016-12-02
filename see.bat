FOR /F "tokens=4 delims= " %%P IN ('netstat -a -n -o ^| findstr :8081') DO @ECHO TaskKill.exe /PID %%P
