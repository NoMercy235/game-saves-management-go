package main

import "net"
import "fmt"
import "bufio"
import (
	//"strings"
	"time"
	//"math/rand"
	//"log"
	"os"
)

var tokenChan chan string = make(chan string)
var token string = "adhuirg38rewbfahd"

func listen (listenPort string)  {
	println("started listening on " + listenPort)
	ln, _ := net.Listen("tcp", "localhost:" + listenPort)
	conn, err := ln.Accept()
	if err != nil {
		println(err)
		return
	}
	for {
		token, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			conn.Close()
			fmt.Println(err)
			break
		}

		fmt.Print("I have the token!\n")
		tokenChan <- token
		conn.Write([]byte("\n"))
	}
}

var test int;
func send (sendPort string) {
	println("started sending on " + sendPort)
	conn, err := net.Dial("tcp", "localhost:" + sendPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	if test == 0 && sendPort == "8081"{
		test = 1
		fmt.Fprintf(conn, token + "\n")
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	for token := range tokenChan {
		time.Sleep(time.Duration(3000) * time.Millisecond)

		fmt.Fprintf(conn, token + "\n")
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func main() {
	test = 0;
	if len(os.Args) != 3 {
		println("Wrong usage")
		return
	}

	listenPort := os.Args[1]
	sendPort := os.Args[2]

	go listen(listenPort)
	go send(sendPort)

	fmt.Scanln()
}