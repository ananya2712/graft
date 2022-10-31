package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
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
	leaderChan       chan bool
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
	lock             sync.Mutex
}

func NodeConstructor(nodeId int, heartBeat int, peerList []int) Node {
	var node Node
	node.state = Follower
	node.nodeId = nodeId
	node.currLeader = -1
	node.currTerm = 0
	node.heartbeat = false
	node.heartbeatTimeout = heartBeat
	node.running = true
	node.peerList = peerList
	node.numServers = len(peerList)
	node.voteFor = -1
	node.voteCount = 0
	node.leaderChan = make(chan bool)
	return node
}

func (node *Node) CloseServer(args *int, reply *int) error {
	node.running = false
	fmt.Println("SERVER CLOSED ", node.nodeId)
	return nil
}

func (node *Node) runRPCListener() {
	srv := rpc.NewServer()
	srv.Register(node)
	//fmt.Println(":" + strconv.Itoa(node.peerList[node.nodeId-1]))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(node.peerList[node.nodeId-1]))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for node.running {
			//fmt.Println("listening", node.nodeId)
			srv.Accept(listener)
		}
	}()
}

func (node *Node) runStateMachine() {
	lastSetTime := time.Now()
	//fmt.Println(node.peerList)
	node.runRPCListener()
	for node.running {
		//fmt.Println(node.nodeId, " node in new iter")
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
			node.voteFor = node.nodeId
			node.voteCount = 1
			go node.broadcastVoteRequest()

			select {
			case <-time.After(4 * time.Second):
				node.state = Follower
			case <-node.leaderChan:
				fmt.Println("leader channel recieved ", node.nodeId)
				node.currTerm++
				fmt.Printf("Node %d is elected leader\n", node.nodeId)
				node.state = Leader
			}

		}
		if node.state == Leader {
			if time.Since(lastSetTime) >= time.Second*time.Duration(node.heartbeatTimeout/2) {
				lastSetTime = time.Now()
				node.broadcastHeartBeat()
			}
		}
	}
}
