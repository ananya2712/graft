package main

import (
	"log"
	"net/rpc"
	"strconv"
)

type heartBeatArgs struct {
	id         int
	currLeader int
}

type heartBeatReply struct {
	success   bool
	id        int
	nextIndex int
}

func (node *Node) broadcastHeartBeat() {
	var args = heartBeatArgs{
		id:         node.nodeId,
		currLeader: node.currLeader,
	}
	_ = args

	for i := range node.peerList {
		go func(i int) {
			var reply heartBeatReply
			node.sendHeartBeat(i, args, &reply)
		}(i)
	}
}

func (node *Node) sendHeartBeat(port int, args heartBeatArgs, reply *heartBeatReply) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
	client.Call("Node.HeartBeat", args, reply)

	if reply.success {
		node.heartbeat = true
	} else {
		node.heartbeat = false
		node.state = Follower
	}
	//incomplete - need to handle nextIndex
}
