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

	client.Call("Node.requestVoteRes", args, reply)

	if reply.currTerm > args.term {
		node.currTerm = reply.currTerm
		node.state = Candidate
		node.voteFor = -1
		return
	}

	if reply.granted {
		node.voteCount++
	}

	if node.voteCount >= len(node.peerList)/2+1 {
		// incomplete - need to make channel for leader announcement
	}
}

func (node *Node) requestVoteRes(args VoteReqArgs, reply *VoteReply) error {

	if args.term < node.currTerm {
		reply.currTerm = node.currTerm
		reply.granted = false
		return nil
	}

	if node.voteFor == -1 {
		node.currTerm = args.term
		node.voteFor = args.candidateId
		reply.currTerm = node.currTerm
		reply.granted = true
	}

	return nil
}
