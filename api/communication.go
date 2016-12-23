package api

import (
	"net"
	"fmt"
	"bufio"
)

// attempts to send a message on a port
// onSuccess: return 0
// onFail: return -1
func Send(self State) (int) {
	// send to self.sendPort
	println("started sending on " + self.SendPort)
	conn, err := net.Dial("tcp", "localhost:" + self.SendPort)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	_, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}


// continuously listening to messages
// onFail: return -1
func Listen(self State) (int) {
	// listen to self.listenPort
	println("started listening on " + self.ListenPort)
	ln, _ := net.Listen("tcp", "localhost:" + self.ListenPort)
	conn, err := ln.Accept()
	if err != nil {
		println(err)
		return -1
	}
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			conn.Close()
			fmt.Println(err)
			return -1
		}
		if msg == "close\n" {
			conn.Close()
			break;
		}

		conn.Write([]byte("\n"))
	}
	return 0
}
