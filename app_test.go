package main

import (
	"testing"
	"os/exec"
	//"time"
)

/*
IMPORTANT:
If this is too hard to implement, we might as well just stick to the uni testing
and leave it be.
 */

/*
This is the integration test for the app. It should test the entire functionality
 */
func TestApp(t *testing.T) {

	if err := exec.Command("start.bat").Run(); err != nil {
		t.Log("Error starting the application")
		t.Log(err)
		t.Fail()
	}
	// Find a way to control the action of a single process
	// For instance, it should shut itself down after some time, if it is
	// a leader, and then check to see if a new leader is re-elected

	//time.Sleep(5 * time.Second)
	//if err := exec.Command("kill.bat").Run(); err != nil {
	//	t.Log("Error closing the application")
	//	t.Log(err)
	//	t.Fail()
	//}
}