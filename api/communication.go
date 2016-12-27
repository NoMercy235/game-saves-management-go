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

	println("mesage sent: " + message)
	fmt.Fprintf(self.SendConn, message + "\n")
	_, err := bufio.NewReader(self.SendConn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		time.Sleep(10000 * time.Millisecond)
		return -1
	}
	return 0
}

// attempts to send a message on a port
// onSuccess: return 0
// onFail: return -1
func SendLeader(self *State, message string) (int) {
	// send to self.sendPort
	println(self.ListenPort + " sends to " + self.LeaderPort)
	if self.SendLeaderConn == nil {
		var err error
		self.SendLeaderConn, err = net.Dial("tcp", "localhost:" + self.LeaderPort)
		if err != nil {
			fmt.Println(err)
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
	}
	for {
		println("INEP SA TRIMIT CATRE LEADER")
		//println("mesage sent: " + message)
		fmt.Fprintf(self.SendLeaderConn, message + "\n")
		println("am trimis de pe " + self.ListenPort)
		_, err := bufio.NewReader(self.SendLeaderConn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
	}
	return 0
}

// continuously listening to messages
// onFail: return -1
func Listen(self *State) (int) {
	// listen to self.listenPort
	println("started listening on " + self.ListenPort + "\n")
	ln, err := net.Listen("tcp", "localhost:" + self.ListenPort)
	if err != nil {
		fmt.Println(err)
		time.Sleep(10000 * time.Millisecond)
		return -1
	}
	if self.ListenConn == nil {
		var err error
		self.ListenConn, err =ln.Accept()
		if err != nil {
			fmt.Println(err)
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
	}

	// this might be helpful in the future to send a timeout. if used, care for the delay below. it should be greater than that
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		msg, err := bufio.NewReader(self.ListenConn).ReadString('\n')
		time.Sleep(500 * time.Millisecond)  // This delay is just to help visualize the communication, since there are a lot of messages.
		if err != nil {
			self.ListenConn.Close()
			fmt.Println(err)
			println("AM DAT EROATE AICI")
			time.Sleep(10000 * time.Millisecond)
			return -1
		}
		println("Received:" + msg[:(len(msg)-1)])
		self.ListenConn.Write([]byte("\n"))
		if msg == "close\n" {
			self.ListenConn.Close()
			break;
		}
		for i := 0; i < len(self.Callbacks); i++ {
			go self.Callbacks[i](self, msg[:(len(msg)-1)]) // msg[:(len(msg)-2)]) it's the msg without the '\n' from the end
		}
	}
	println("*** Process " + self.ListenPort + " is no longer listening! ***")
	return 0
}
//
//func startListeningForConnections(ln net.Listener) (int) {
//	println("Leader TRYING TO LISTEN ON --- " + currentProcess.ListenPort)
//	conn, _ := ln.Accept()
//	for {
//		dec := gob.NewDecoder(conn)
//		message := &LeaderElectionMessage{}
//		dec.Decode(message)
//		if (message.IsPing == true) {
//			println("Recieved ping request");
//		}
//		conn.Write([]byte("\n"));
//	}
//	return 0
//}
