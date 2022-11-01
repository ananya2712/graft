package main

import (
	"log"
	"os"
	"strconv"
	"sync"
)

var cleader CurrentLeader

func createCluster(n int, heartBeat int) {
	var clusterwg sync.WaitGroup
	var startPort int = 8080
	var peerList []int
	cleader = CurrentLeaderConstructor(1, 0)

	for i := 1; i < n+1; i++ {
		peerList = append(peerList, startPort+i)
	}

	for i := 0; i < n; i++ {
		tempNode := NodeConstructor(i+1, heartBeat, peerList)
		go tempNode.runStateMachine()
		if i == 0 {
			tempNode.state = Leader
		}
		clusterwg.Add(1)
	}

	clusterwg.Wait()
}

func main() {

	numNodes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	heartBeat, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	createCluster(numNodes, heartBeat)
}
