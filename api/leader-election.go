package api

import (
	"fmt"

	"net"
	"bufio"

	"encoding/gob"
	"time"
)

//import "time"


/*
This file should do the following:

 If i am a leader:
 - respond to others that are trying to see if I'm still alive

 If i am not a leader:
 - check to see if the leader still lives, and if it is not, initiate the leader election algorithm


 Leader election algorithm:
 We will use the token technique, visualizing the topology in a ring. The process that has noticed that the leader
 is gone, will initiate this algorithm. It will instantiate a random token (we should make a struct for this), and will
 fill a property of it with it PID then pass it to the next in circle (state.SendPort).
 The token will get passed like that and every process will fill the PID property with their PID *only* if their PID
 is grater than what was already there.
 When a process gets the token and find the same PID as its own, it then becomes a leader and sends a message to all the
 others to notify them (they all set the state.LeaderPort property

 IMPORTANT: when starting the app, all the processes will have "" in the LeaderPort. Which means that all of them will
 initiate the algorithm, leading to many tokens travelling through the topology. Should use a lot of logs to see how
 exactly all of this is happening
 */

var currentProcess State;

var tokenChan chan *LeaderElectionMessage = make(chan *LeaderElectionMessage);

var listener net.Listener;

func CheckLeader(self State) {

	currentProcess = self;
	go listenForLeader()
	go sendForLeader(nil)

}

func sendForLeader(conn net.Conn) (int) {
	// send to self.sendPort
	for len(tokenChan) > 0 {
		<-tokenChan
	}
	println("started sending on " + currentProcess.SendPort)
	var err error;
	conn, err = net.Dial("tcp", "localhost:" + currentProcess.SendPort)
	if (err != nil) {
		fmt.Println(err)
		return -1;
	}

	encoder := gob.NewEncoder(conn)
	p := &LeaderElectionMessage{}
	p.LeaderPort = currentProcess.LeaderPort;
	encoder.Encode(p)
	println("Sent first message");

	for message := range tokenChan {
		x := false
		time.Sleep(time.Duration(2000) * time.Millisecond)
		encoder := gob.NewEncoder(conn)
		encoder.Encode(message)

		small:
		if (x == false) {
			_, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				println("66666")
				fmt.Println(err)
				tokenChan <- message
				x = true
				goto small
			}
		}
		println("Sent message on port " + currentProcess.SendPort);
	}

	conn.Close()
	return 0
}

func listenForLeader() (int) {
	// listen to self.listenPort
	println("started listening on " + currentProcess.ListenPort)
	if (listener != nil) {
		listener.Close()
	}

	var err error;
	listener, err = net.Listen("tcp", "localhost:" + currentProcess.ListenPort)
	if (err != nil) {
		fmt.Println(err)
		return -1;
	}
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	for {
		println("Waiting for a message on port : " + currentProcess.ListenPort)
		dec := gob.NewDecoder(conn)
		message := &LeaderElectionMessage{}
		dec.Decode(message)

		println("------------------");
		fmt.Printf("Received 1 : %+v", message);
		println("");
		println("------------------");

		println("My PID == " + currentProcess.ListenPort);
		println("Received PID == " + message.LeaderPort);

		if (message.LeaderPort < currentProcess.ListenPort && message.LeaderFound == false) {
			println("My PID is higher, setting me as the leader");
			message.LeaderPort = currentProcess.ListenPort;
			tokenChan <- message;
			conn.Write([]byte("\n"))
		} else if (message.LeaderPort > currentProcess.ListenPort && message.LeaderFound == false) {
			println("My PID is lower, not modifying the message");
			tokenChan <- message;
			conn.Write([]byte("\n"))
		} else if (message.LeaderPort == currentProcess.ListenPort ) {
			println("I AM THE LEADER! OBEY ME, PEASANTS");
			currentProcess.IsLeader = true;
			currentProcess.LeaderPort = currentProcess.ListenPort;
			message.LeaderFound = true;
			message.LeaderPort = currentProcess.ListenPort;
			message.LeaderSendPort = currentProcess.SendPort
			println("yes loop")
			tokenChan <- message
			conn.Write([]byte("\n"))
			go startListeningForConnections(listener);
			go startListeningForConnections(listener);
			message.FirstLoop = true;
		} else if (currentProcess.IsLeader == false && message.LeaderFound == true) {
			println("Leader found (" + currentProcess.LeaderPort +")" + "I will ping him from now on")
			currentProcess.LeaderPort = message.LeaderPort;
			currentProcess.LeaderSendPort = message.LeaderSendPort;
			tokenChan <- message;
			conn.Write([]byte("\n"));
			if (currentProcess.ListenPort != "80281") {
				conn.Write([]byte("\n"));
				go startPingingLeader()
				break;
			}
		}
		println("");
		fmt.Printf("Current Process 2 : %+v", currentProcess);
		println("");
	}

	return 0
}

func startListeningForConnections(ln net.Listener) (int) {
	println("Leader TRYING TO LISTEN ON --- " + currentProcess.ListenPort)

	conn, _ := ln.Accept()
	for {
		dec := gob.NewDecoder(conn)
		message := &LeaderElectionMessage{}
		dec.Decode(message)
		if (message.IsPing == true) {
			println("Recieved ping request");
		}

		conn.Write([]byte("\n"));
	}

	return 0
}

func startPingingLeader() (int) {
	// send to self.sendPort
	time.Sleep(time.Duration(2222) * time.Millisecond)
	println("My PID == " + currentProcess.ListenPort);
	println("Started pinging leader on port : " + currentProcess.LeaderPort)
	conn, err := net.Dial("tcp", "localhost:" + currentProcess.LeaderPort)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		println("1")
		encoder := gob.NewEncoder(conn)
		p := &LeaderElectionMessage{}
		p.IsPing = true;
		encoder.Encode(p)

		println("Sent ping to leader on port " + currentProcess.LeaderPort);
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			println("WE NEED A NEW LEADER");
			removeLeader(currentProcess.AllPorts)
			currentProcess.IsLeader = false
			if (currentProcess.SendPort == currentProcess.LeaderPort) {
				currentProcess.SendPort = currentProcess.LeaderSendPort;
			}

			currentProcess.LeaderPort = ""
			fmt.Printf("%v", currentProcess)
			go CheckLeader(currentProcess)
			break;
		}

	}

	return 0
}
func removeLeader(allPorts []string) {
	for i := 0; i < len(allPorts); i++ {
		if (currentProcess.LeaderPort == allPorts[i]) {
			allPorts = append(allPorts[:i], allPorts[i + 1:]...)
			break;
		}
	}

	currentProcess.AllPorts = allPorts;
}
