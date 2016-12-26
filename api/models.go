package api

type State struct {
	ListenPort string
	SendPort string
	AllPorts []string
	LeaderPort string
	IsLeader bool
	Callbacks []func(self *State, message string)
}

func RegisterCallback(self *State, function func(self *State, message string)) {
	println("am pus callback:  " + GetFunctionName(function))
	self.Callbacks = append(self.Callbacks, function)
}