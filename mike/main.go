package main

import (
	"context"
	"fmt"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func lookup (nodeAddr, key string) (*pb.Response, error) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 8000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	request := &pb.CheckFileRequest{
		Key:     key,
		Command: pb.Command_CHECK,
	}
	client := pb.NewExternalListenerClient(conn)
	resp, err := client.CheckFile(context.Background(), request)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return nil, err
	}
	log.Printf("Redirect is %v, for node %v\n", resp.Redirect, resp.NodeInfo.NodeID )
	return resp, nil
}

func main(){
	contactNode := os.Getenv("APP_NODE")

	for {
		time.Sleep(2*time.Second)
		res, err := lookup(contactNode, "ORD")
		fmt.Println("Pingpong zimzam")
		if err != nil {
			log.Fatalf("Crippling depression")
		}
		log.Printf("Result is %+v", res)
	}
}


