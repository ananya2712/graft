package main

import (
	"fmt"
	"time"
)

const (
	Follower  int = 0
	Candidate int = 1
	Leader    int = 2
)

type Node struct {
	state      int
	nodeId     int
	currTerm   int
	currLeader int
	heartbeat  bool
}

func NodeConstructor(nodeId int) Node {
	var node Node
	node.state = Follower
	node.nodeId = nodeId
	node.currLeader = 0
	node.currTerm = 0
	node.heartbeat = false
	return node
}

func (node *Node) runStateMachine() {
	lastSetTime := time.Now()
	for {
		if node.state == Follower {
			if time.Since(lastSetTime) >= time.Second*3 {
				fmt.Println("inside")
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
