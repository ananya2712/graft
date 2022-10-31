package main

import (
	"log"
	"net/rpc"
	"strconv"
)

type heartBeatArgs struct {
	CurrLeader int
	CurrTerm   int
}

type heartBeatReply struct {
	Success bool
	NodeId  int
}

func (node *Node) broadcastHeartBeat() {
	var args = heartBeatArgs{
		CurrLeader: node.nodeId,
		CurrTerm:   node.currTerm,
	}
	_ = args

	for i := range node.peerList {
		go func(i int) {
			var reply heartBeatReply
			node.sendHeartBeat(i, args, &reply)
		}(i)
	}
}

func (node *Node) HeartBeat(args *heartBeatArgs, reply *heartBeatReply) error {
	reply.NodeId = node.nodeId
	if args.CurrLeader == node.currLeader && args.CurrTerm == node.currTerm {
		node.heartbeat = true
		reply.Success = true
	} else {
		node.heartbeat = false
		reply.Success = false
	}
	return nil
}

func (node *Node) sendHeartBeat(port int, args heartBeatArgs, reply *heartBeatReply) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
	client.Call("Node.HeartBeat", args, reply)
}
