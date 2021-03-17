package main

import (
	"context"
	"fmt"
	gossip "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func lookup(nodeAddr, key string) {
	log.Printf("Attempting to lookup string %v", key)
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	checkRequest := &gossip.CheckRequest{Key: key}
	client := gossip.NewInternalListenerClient(conn)
	resp, err := client.CheckIP(context.Background(), checkRequest)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	log.Printf("Got the result for key %v found in node %v\n", key, resp.IP)
}

func main(){
	contactNode := os.Getenv("APP_NODE")
	keys := [4]string{"AAA", "BBB", "XXX", "UYEWTFBBQ"}

	for i:=0; i < 999; i++{
		fmt.Println("Trying a new key...")
		// harcoded AAA -> bravo BBB -> alpha XXX -> charlie
		lookup(contactNode, keys[i % 4])
		time.Sleep(2*time.Second)


	}
}


