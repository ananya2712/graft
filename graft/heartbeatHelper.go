package main

import (
	"log"
	"net/rpc"
	"strconv"
)

type HeartBeatArgs struct {
	CurrLeader int
	CurrTerm   int
}

type HeartBeatReply struct {
	Success bool
	NodeId  int
}

func (node *Node) broadcastHeartBeat() {
	var args = HeartBeatArgs{
		CurrLeader: node.nodeId,
		CurrTerm:   node.currTerm,
	}
	_ = args

	for i := range node.peerList {
		go func(i int) {
			var reply HeartBeatReply
			node.sendHeartBeat(node.peerList[i], args, &reply)
		}(i)
	}
}

func (node *Node) HeartBeat(args *HeartBeatArgs, reply *HeartBeatReply) error {
	//fmt.Println("Heartbeat receieved!", node.nodeId, " from ", args.CurrLeader, " ", args.CurrTerm)
	node.heartbeat <- true
	node.currTerm = args.CurrTerm
	return nil
}

func (node *Node) sendHeartBeat(port int, args HeartBeatArgs, reply *HeartBeatReply) {
	client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
	//fmt.Println("Sending heartbeat from ", node.nodeId)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
	err = client.Call("Node.HeartBeat", args, reply)
	if err != nil {
		log.Fatal(err)
	}
}
