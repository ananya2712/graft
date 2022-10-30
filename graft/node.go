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
	state     int
	nodeId    int
	term      int
	leader    int
	heartbeat bool
}

func NodeConstructor(nodeId int) Node {
	var node Node
	node.state = Follower
	node.nodeId = nodeId
	return node
}

func (node *Node) runStateMachine() {
	lastSetTime := time.Now()
	for {
		if node.state == Follower {
			if time.Since(lastSetTime) >= time.Second*3 && !node.heartbeat {
				node.state = Candidate
				fmt.Println(node.nodeId, " is a candidate")
			} else {
				node.heartbeat = !node.heartbeat
				fmt.Println(node.nodeId, " has receieved a heartbeat from leader")
			}
			lastSetTime = time.Now()
		}
		if node.state == Candidate {

		}
		if node.state == Leader {

		}
	}
}
