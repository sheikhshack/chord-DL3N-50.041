package main

import (
	"context"
	"fmt"
	exposed "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func lookup (nodeAddr, key string) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 8000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	checkRequest := &exposed.Request{Key: key}
	client := exposed.NewExternalListenerClient(conn)
	resp, err := client.CheckIP(context.Background(), checkRequest)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	log.Printf("Got the result for key %v : %+v\n", resp.IP)
}

func main(){
	contactNode := os.Getenv("APP_NODE")

	for {
		time.Sleep(2*time.Second)
		lookup(contactNode, "ORD")

	}
}


