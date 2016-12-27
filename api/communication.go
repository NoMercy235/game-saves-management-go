package api

/*
This file will be used only for communication. Basically, the functions Send and Listen should do all that is required
for now but they may be subject to change
 */

import (
	"net"
	"fmt"
	"bufio"
	"time"
)


// attempts to send a message on a port
// onSuccess: return 0
// onFail: return -1
func Send(self *State, message string) (int) {
	// send to self.sendPort
	println(self.ListenPort + " sends to " + self.SendPort)
	if self.SendConn == nil {
		var err error
		self.SendConn, err = net.Dial("tcp", "localhost:" + self.SendPort)
		if err != nil {
			fmt.Println(err)
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
	}

	fmt.Fprintf(self.SendConn, message + "\n")
	_, err := bufio.NewReader(self.SendConn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		time.Sleep(10000 * time.Millisecond)
		return -1
	}
	return 0
}
//var test int;
//func Send(self *State, message string) (int) {
//	println(self.ListenPort + " sends to " + self.SendPort)
//	conn, _ := net.Dial("tcp", "localhost:" + self.SendPort)
//
//	if test == 0 && self.ListenPort == "8081" {
//		test = 1
//		fmt.Fprintf(conn, message + "\n")
//		_, err := bufio.NewReader(conn).ReadString('\n')
//		if err != nil {
//			fmt.Println(err)
//			return 0
//		}
//	}
//
//	for token := range self.TokenChan {
//		fmt.Fprintf(conn, token + "\n")
//		_, err := bufio.NewReader(conn).ReadString('\n')
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//	}
//	return 0
//}


// continuously listening to messages
// onFail: return -1
func Listen(self *State) (int) {
	// listen to self.listenPort
	println("started listening on " + self.ListenPort + "\n")
	ln, _ := net.Listen("tcp", "localhost:" + self.ListenPort)
	conn, err := ln.Accept()
	if err != nil {
		println(err)
		return -1
	}
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		time.Sleep(1000 * time.Millisecond)
		if err != nil {
			conn.Close()
			fmt.Println(err)
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
		println("Received:" + msg)
		conn.Write([]byte("\n"))
		if msg == "close\n" {
			conn.Close()
			break;
		}
		//self.TokenChan <- msg
		for i := 0; i < len(self.Callbacks); i++ {
			go self.Callbacks[i](self, msg[:(len(msg)-2)]) // msg[:(len(msg)-2)]) it's the msg without the '\n' from the end
		}
	}
	return 0
}