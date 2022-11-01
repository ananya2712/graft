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
	//CurrTerm int
	Granted bool
}

type NewLeaderArgs struct {
	LeaderId int
	Term     int
}

func (node *Node) NewLeader(args NewLeaderArgs, reply *int) error {
	node.state = Follower
	node.currTerm = args.Term
	return nil
}

func (node *Node) RequestVoteRes(args VoteReqArgs, reply *VoteReply) error {
	fmt.Println(node.nodeId, args.NodeId, args.Term, node.currTerm, node.voteFor)
	if args.Term < node.currTerm {
		//fmt.Println("Inside cond 1")
		reply.Granted = false
		return nil
	}

	if node.voteFor == -1 {
		//fmt.Println("Inside cond 2")
		reply.Granted = true
		node.lock.Lock()
		node.voteFor = args.NodeId
		node.lock.Unlock()
	}

	return nil
}

func (node *Node) sendVoteReq(port int, args VoteReqArgs, reply *VoteReply) {
	client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
	err = client.Call("Node.RequestVoteRes", args, reply)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(node.nodeId, node.state)
	fmt.Println(reply)

	if reply.Granted {
		node.voteCount++
		fmt.Println(node.nodeId, " votes:", node.voteCount)
		if node.voteCount >= len(node.peerList)/2+1 {
			//fmt.Println("Leader possible ", node.nodeId)
			node.voteCount = 0
			node.voteFor = -1
			node.leaderChan <- true
		}
	}
}

func (node *Node) broadcastVoteRequest() {
	var args = VoteReqArgs{
		Term:   node.currTerm + 1,
		NodeId: node.nodeId,
	}

	for i := range node.peerList {
		if true {
			go func(i int) {
				var reply VoteReply
				node.sendVoteReq(node.peerList[i], args, &reply)
			}(i)
		}
	}
}

func (node *Node) broadcastNewLeader() {
	var args = NewLeaderArgs{
		Term:     node.currTerm,
		LeaderId: node.nodeId,
	}
	for i := range node.peerList {
		go func(i int) {
			var reply int
			client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(node.peerList[i]))
			if err != nil {
				log.Fatal("dialing:", err)
			}
			defer client.Close()
			err = client.Call("Node.NewLeader", args, &reply)
			if err != nil {
				log.Fatalln(err)
			}
		}(i)
	}
}
