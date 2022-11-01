package main

import (
	"log"
	"net/rpc"
	"strconv"
)

type LogRecord struct {
	LogId    int
	LogValue int
}

type AddValueArgs struct {
	Value int
}

func (node *Node) SendValue(args AddValueArgs, reply *bool) error {
	if len(node.log) > 0 {
		oldLog := node.log[len(node.log)-1]
		newLog := LogRecord{oldLog.LogId + 1, args.Value}
		node.log = append(node.log, newLog)
	} else {
		newLog := LogRecord{0, args.Value}
		node.log = append(node.log, newLog)
	}
	return nil
}

func (node *Node) AddValue(args AddValueArgs, reply *bool) error {
	log.Println("Node", node.nodeId, "has received value update request")
	if node.state == Leader {
		if len(node.log) > 0 {
			oldLog := node.log[len(node.log)-1]
			newLog := LogRecord{oldLog.LogId + 1, args.Value}
			node.log = append(node.log, newLog)
		} else {
			newLog := LogRecord{0, args.Value}
			node.log = append(node.log, newLog)
		}
		for i := 0; i < len(node.peerList); i++ {
			if i != node.nodeId-1 {
				go func(i int) {
					client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(node.peerList[i]))
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()
					var reply bool
					err = client.Call("Node.SendValue", args, &reply)
					if err != nil {
						log.Fatal(err)
					}
				}(i)
			}
		}
		*reply = true
	} else {
		log.Println("Node", node.nodeId, "is not a leader. Failure.")
		*reply = false
	}
	return nil
}

func (node *Node) GetValue(args int, reply *LogRecord) error {
	log.Println("Node", node.nodeId, "has received value read request")
	if len(node.log) > 0 {
		*reply = node.log[len(node.log)-1]
	} else {
		*reply = LogRecord{-1, -1}
	}
	return nil
}
