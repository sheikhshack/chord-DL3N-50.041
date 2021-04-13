package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gossip "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"

	"google.golang.org/grpc"
)

func writeExternalFile(nodeAddr, fileName, containerIP string) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.StoreKeyHash(context.Background(), &gossip.DLUploadRequest{
		Filename:    fileName,
		ContainerIP: containerIP,
	})
	if err != nil {
		log.Fatalf("[MIKE] external fail sending to %v: %s", nodeAddr, err)
	}

	log.Printf("\nSuccess upload info to the following chord node: %+v\n", res)
}

func resolveFile(nodeAddr, fileName string) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.GetFileLocation(context.Background(), &gossip.DLDownloadRequest{
		Filename: fileName,
	})
	if err != nil {
		log.Fatalf("[MIKE] external fail sending to %v: %s", nodeAddr, err)
	}

	log.Printf("\nSuccess, received the following containerINFO: %+v\n", res)
}

func main() {
	fmt.Printf("%.23s\n", time.Now().UTC())
	attachedNode := os.Getenv("APP_NODE")
	fileName := os.Getenv("FILE_NAME")
	containerIP := os.Getenv("CONTAINER_IP")

	fileNames := strings.Split(fileName, ",")
	for _, i := range fileNames {
		//fmt.Println("Sleeping the night away")
		//time.Sleep(5 * time.Second)
		writeExternalFile(attachedNode, i, containerIP)
		// fmt.Println("Sleep again, retrieving back the same file in 5s")
		// time.Sleep(5 * time.Second)
		// resolveFile(attachedNode, i)
	}

}
