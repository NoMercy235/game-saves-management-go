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
func Send(self *State, port string, message string, print bool) (int) {
	//println(self.ListenPort + " sends to " + port)
	if self.Connections[port] == nil {
		var err error
		self.Connections[port], err = net.Dial("tcp", "localhost:" + port)
		if err != nil {
			fmt.Println(err)
			time.Sleep(ERR_TIME)
			return -1
		}
	}

	if print {
		println("message sent to " + port + ": " + message)
	}
	fmt.Fprintf(self.Connections[port], message + "\n")
	_, err := bufio.NewReader(self.Connections[port]).ReadString('\n')
	if err != nil {
		self.Connections[port].Close()
		fmt.Println(err)
		time.Sleep(ERR_TIME)
		return -1
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
		time.Sleep(ERR_TIME)
		return -1
	}

	// this might be helpful in the future to send a timeout. if used, care for the delay below. it should be greater than that
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		println("*** Accepted a connection! ***")
		conn, _ := ln.Accept()
		go handleConnection(self, conn)
	}
	println("*** Process " + self.ListenPort + " is no longer listening! ***")
	return 0
}

func handleConnection(self *State, conn net.Conn){
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		time.Sleep(DELAY_TIME)  // This delay is just to help visualize the communication, since there are a lot of messages.
		if err != nil {
			conn.Close()
			fmt.Println(err)
			time.Sleep(ERR_TIME)
			break
		}
		println("Received:" + msg[:(len(msg)-1)])
		conn.Write([]byte("\n"))
		if msg == "close\n" {
			conn.Close()
			break;
		}
		for i := 0; i < len(self.Callbacks); i++ {
			go self.Callbacks[i](self, msg[:(len(msg)-1)]) // msg[:(len(msg)-2)]) it's the msg without the '\n' from the end
		}
	}
}




/*
Those two function only serve as a connectivity test (for now) between all the processes of the application
 */
func PingEveryone(self *State) {
	for i := 0; i < len(self.AllPorts); i ++ {
		if self.AllPorts[i] != self.ListenPort {
			go ping(self, self.AllPorts[i])
		}
	}
}

func ping(self *State, port string) {
	println("SENDING TO " + port)
	conn, err := net.Dial("tcp", "localhost:" + port)
	if err != nil {
		fmt.Println(err)
		time.Sleep(ERR_TIME)
		return
	}

	for {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(conn, "PING FROM" + self.ListenPort + "\n")
		_, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			time.Sleep(ERR_TIME)
			return
		}
	}
	return
}
