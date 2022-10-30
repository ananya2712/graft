package main

import (
	"log"
	"os"
	"strconv"
	"sync"
)

func createCluster(n int) {
	var clusterwg sync.WaitGroup
	for i := 0; i < n; i++ {
		newNode := NodeConstructor(i + 1)
		go newNode.runStateMachine()
		clusterwg.Add(1)
	}
	clusterwg.Wait()
}

func main() {
	numNodes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	createCluster(numNodes)
}
