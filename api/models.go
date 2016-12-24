package api

type State struct {
	ListenPort string
	SendPort string
	AllPorts []string
	LeaderPort string
	IsLeader bool
}