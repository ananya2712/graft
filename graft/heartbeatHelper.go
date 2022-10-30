package main

type heartBeatArgs struct {
	id         int
	currLeader int
}

type heartBeatReply struct {
	success bool
	id      int
}
