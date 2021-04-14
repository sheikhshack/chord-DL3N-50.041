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

func deleteFile(nodeAddr, fileName string) {
	//log.Printf("Attempting to lookup string %v", key)
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.DeleteClientFile(context.Background(), &gossip.DLDeleteRequest{
		Filename: fileName,
	})
	if err != nil {
		log.Fatalf("-- MIKE external fail %s", err)
	}

	log.Printf("\nSuccess, deleted the following containerINFO: %+v\n", res)
}

func main() {
	fmt.Printf("%.23s\n", time.Now().UTC())
	attachedNode := os.Getenv("APP_NODE")
	fileName := os.Getenv("FILE_NAME")
	containerIP := os.Getenv("CONTAINER_IP")

	fileNames := strings.Split(fileName, ",")

	// File Upload test
	for _, i := range fileNames {
		//fmt.Println("Sleeping the night away")
		//time.Sleep(5 * time.Second)
		writeExternalFile(attachedNode, i, containerIP)
	}

	// // Faulty - File lookup test
	// fmt.Println("Sleep again for 20s")
	// time.Sleep(20 * time.Second)

	// for _, i := range fileNames {
	// 	fmt.Println("Sleep again, retrieving back the same file in 5s")
	// 	time.Sleep(5 * time.Second)
	// 	resolveFile("bravo", i)
	// }

	// File delete test
	fmt.Println("Sleep again for 10s")
	time.Sleep(10 * time.Second)

	for _, i := range fileNames {
		fmt.Println("Sleep again, deleting the same file in 5s")
		time.Sleep(5 * time.Second)
		deleteFile(attachedNode, i)
	}
}
