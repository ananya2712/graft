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

type CurrentLeader struct {
	LeaderId   int
	Term       int
	LeaderLock sync.Mutex
}

type RPCArguments struct {
	nodeId bool
	hello  string
}

type Node struct {
	leaderChan       chan bool
	state            int
	peerList         []int
	nodeId           int
	currTerm         int
	currLeader       *CurrentLeader
	heartbeat        chan bool
	heartbeatTimeout int
	running          bool
	numServers       int
	log              []LogRecord
	voteFor          int
	voteCount        int
	lock             sync.Mutex
}

func CurrentLeaderConstructor(nodeId int, term int) CurrentLeader {
	var currentLeader CurrentLeader
	currentLeader.LeaderId = nodeId
	currentLeader.Term = term
	currentLeader.LeaderLock = sync.Mutex{}
	return currentLeader
}

func NodeConstructor(nodeId int, heartBeat int, peerList []int) Node {
	var node Node
	node.state = Follower
	node.nodeId = nodeId
	//node.currLeader = -1
	node.currTerm = 0
	node.heartbeat = make(chan bool)
	node.heartbeatTimeout = heartBeat
	node.running = true
	node.peerList = peerList
	node.numServers = len(peerList)
	node.voteFor = -1
	node.voteCount = 0
	node.leaderChan = make(chan bool)
	return node
}

func (node *Node) CloseServer(args int, reply *int) error {
	node.running = false
	log.Println("Server Shut Down", node.nodeId)
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
			select {
			case <-time.After(time.Second * time.Duration(node.heartbeatTimeout)):
				log.Println("Did not recieve heartbeat ", node.nodeId)
				//fmt.Println(node.nodeId, node.currLeader)
				cleader.LeaderLock.Lock()
				//fmt.Println("hello from the inside", node.nodeId, cleader.Term, cleader.LeaderId, node.currTerm)
				if cleader.Term <= node.currTerm {
					cleader.LeaderId = node.nodeId
					cleader.Term = node.currTerm + 1
					node.state = Leader
					node.currTerm += 1
				}
				cleader.LeaderLock.Unlock()
			case <-node.heartbeat:
				log.Println("Received heartbeat ", node.nodeId)
				continue
			}
		}
		if node.state == Candidate {

			log.Println(node.nodeId, " is requesting votes")
			//node.voteFor = node.nodeId
			node.voteCount = 1
			go node.broadcastVoteRequest()

			select {
			case <-time.After(time.Second * time.Duration(node.heartbeatTimeout)):
				node.lock.Lock()
				node.state = Follower
				node.lock.Unlock()
			case <-node.leaderChan:
				//fmt.Println("leader channel recieved ", node.nodeId)
				if node.state == Candidate {
					node.currTerm++
					log.Printf("Node %d is elected leader\n", node.nodeId)
					node.lock.Lock()
					node.state = Leader
					fmt.Println(node.voteCount)
					node.lock.Unlock()
					node.broadcastNewLeader()
				}
			}

		}
		if node.state == Leader {
			if time.Since(lastSetTime) >= time.Second*time.Duration(node.heartbeatTimeout/2) {
				log.Println("Leader:", node.nodeId)
				lastSetTime = time.Now()
				node.broadcastHeartBeat()
				//node.broadcastNewLeader()
			}
		}
	}
}
