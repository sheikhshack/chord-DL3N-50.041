package main

import (
	"context"
	"fmt"
	gossip "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

func lookup(nodeAddr, key string) {
	//log.Printf("Attempting to lookup string %v", key)
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
	log.Printf("key %v (hash of %v) found in node %v\n", key, hash.Hash(key), resp.IP)
	upload(resp.IP, key, fmt.Sprintf("%.23s | %s\ndistributed systems is the best :)\n", time.Now().UTC(), key))
}

func upload(nodeAddr, key, value string) {
	//log.Printf("Attempting to lookup string %v", key)
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	req := &gossip.UploadRequest{Key: key, Value: value}
	client := gossip.NewInternalListenerClient(conn)
	_, err = client.Upload(context.Background(), req)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	//log.Printf("key %v uploaded to node %v\n", key, nodeAddr)
}

func main() {
	fmt.Printf("%.23s\n", time.Now().UTC())
	contactNode := os.Getenv("APP_NODE")
	searchKeys := os.Getenv("SEARCH_KEY")
	//keys := [4]string{"AAA", "BBB", "XXX", "UYEWTFBBQ"}

	//for i := 0; i < 999; i++ {
	//	fmt.Println("Trying a new key...")
	//	lookup(contactNode, keys[i%4])
	//	time.Sleep(2 * time.Second)
	//}
	keys := strings.Fields(searchKeys)
	for _, key := range keys {
		lookup(contactNode, key)
	}

}
