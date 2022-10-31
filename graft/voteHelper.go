package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

type VoteReqArgs struct {
	Term   int
	NodeId int
}

type VoteReply struct {
	CurrTerm int
	Granted  bool
}

func (node *Node) RequestVoteRes(args VoteReqArgs, reply *VoteReply) error {
	fmt.Println(node.nodeId, args.NodeId, args.Term, node.currTerm, node.voteFor)
	if args.Term < node.currTerm {
		fmt.Println("Inside cond 1")
		reply.CurrTerm = node.currTerm
		reply.Granted = false
		return nil
	}

	if node.voteFor == -1 {
		fmt.Println("Inside cond 2")
		node.currTerm = args.Term
		node.voteFor = args.NodeId
		reply.CurrTerm = node.currTerm
		reply.Granted = true
	}

	return nil
}

func (node *Node) sendVoteReq(port int, args VoteReqArgs, reply *VoteReply) {

	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	err = client.Call("Node.RequestVoteRes", args, reply)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("reply")
	fmt.Println(reply)

	if reply.Granted {
		node.voteCount++
		if node.voteCount >= len(node.peerList)/2+1 {
			node.leaderChan <- true
			node.voteCount = 0
		}
	} else {
		node.currTerm = reply.CurrTerm
		node.state = Follower
		node.voteFor = -1
		return
	}
}

func (node *Node) broadcastVoteRequest() {
	var args = VoteReqArgs{
		Term:   node.currTerm + 1,
		NodeId: node.nodeId,
	}

	for i := range node.peerList {
		go func(i int) {
			var reply VoteReply
			node.sendVoteReq(node.peerList[i], args, &reply)
		}(i)
	}
}
