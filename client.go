package main

import (
	"flag"
	"fmt"
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

func main() {
	// Stop a server
	// Send a new value to the leader ()
	// Receieve a value from any of the servers
	send_val := flag.Bool("set", false, "send_val")
	get_val := flag.Bool("get", false, "get_val")
	close_server := flag.Bool("close", false, "close_server")
	value := flag.Int("val", 0, "value")
	server := flag.Int("server", 0, "server")
	flag.Parse()

	client, err := rpc.Dial("tcp", "127.0.0.1:808"+strconv.Itoa(*server))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if *close_server {
		a := 0
		err = client.Call("Node.CloseServer", 1, &a)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success: Server", *server, "closed")
		return
	}

	if *get_val {
		res := new(LogRecord)
		err = client.Call("Node.GetValue", 1, res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Current Value in the log:", res.LogValue)
		return
	}
	if *send_val {
		var reply bool
		args := AddValueArgs{*value}
		err = client.Call("Node.AddValue", args, &reply)
		if err != nil {
			log.Fatal(err)
		}
		if !reply {
			fmt.Println("Failure. Send new value to leader")
		} else {
			fmt.Println("Success")
		}
		return
	}
}
