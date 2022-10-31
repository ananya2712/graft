package main

import (
	"log"
	"net/rpc"
	"strconv"
)

type VoteReqArgs struct {
	term        int
	candidateId int
}

type VoteReply struct {
	currTerm int
	granted  bool
}

func (node *Node) broadcastVoteRequest() {
	var args = VoteReqArgs{
		term:        node.currTerm,
		candidateId: node.nodeId,
	}
	_ = args

	for i := range node.peerList {
		go func(i int) {
			var reply VoteReply
			node.sendRequestVote(i, args, &reply)
		}(i)
	}
}

func (node *Node) sendRequestVote(port int, args VoteReqArgs, reply *VoteReply) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
}
