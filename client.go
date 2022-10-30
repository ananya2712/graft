package main

import (
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	a := 0
	err = client.Call("Node.CloseServer", 1, &a)
	if err != nil {
		log.Fatal(err)
	}
}
