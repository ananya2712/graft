package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	Follower  int = 0
	Candidate int = 1
	Leader    int = 2
)

type RPCArguments struct {
	nodeId bool
	hello  string
}

type LogRecord struct {
	logId    int
	logValue int
}

type Node struct {
	state            int
	peerList         []int
	nodeId           int
	currTerm         int
	currLeader       int
	heartbeat        bool
	heartbeatTimeout int
	running          bool
	numServers       int
	log              []LogRecord
	voteFor          int
	voteCount        int
}

func NodeConstructor(nodeId int, heartBeat int, peerList []int) Node {
	var node Node
	node.state = Follower
	node.nodeId = nodeId
	node.currLeader = 0
	node.currTerm = 0
	node.heartbeat = false
	node.heartbeatTimeout = heartBeat
	node.running = true
	node.peerList = peerList
	node.numServers = len(peerList)
	return node
}

func (node *Node) CloseServer(args *int, reply *int) error {
	node.running = false
	fmt.Println("SERVER CLOSED")
	return nil
}

func (node *Node) runRPCListener() {
	rpc.Register(node)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(listener, nil)
}

func (node *Node) runStateMachine() {
	lastSetTime := time.Now()
	node.runRPCListener()
	for node.running {
		if node.state == Follower {
			if time.Since(lastSetTime) >= time.Second*time.Duration(node.heartbeatTimeout) {
				//fmt.Println("inside")
				if !node.heartbeat {
					node.state = Candidate
					fmt.Println(node.nodeId, " is a candidate")
				} else {
					fmt.Println(node.nodeId, " has receieved heartbeat from leader")
					node.heartbeat = false
					lastSetTime = time.Now()
				}
			}
		}
		if node.state == Candidate {
			fmt.Println(node.nodeId, " is requesting votes")
		}
		if node.state == Leader {
			fmt.Println(node.nodeId, " is leader")
		}
	}
}
