package api

type State struct {
	ListenPort string
	SendPort string
	AllPorts []string
	LeaderPort string
	IsLeader bool
	LeaderSendPort string
}

type LeaderElectionMessage struct {
	LeaderPort string
	LeaderFound bool
	LeaderSendPort string
	FirstLoop bool
	IsPing bool
	IsPong bool
}