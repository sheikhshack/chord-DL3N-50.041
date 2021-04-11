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

//func lookup(nodeAddr, key string) {
//	//log.Printf("Attempting to lookup string %v", key)
//	var conn *grpc.ClientConn
//	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)
//
//	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("Cannot connect to: %s", err)
//
//	}
//	defer conn.Close()
//
//	checkRequest := &gossip.CheckRequest{Key: key}
//	client := gossip.NewInternalListenerClient(conn)
//	resp, err := client.CheckIP(context.Background(), checkRequest)
//	if err != nil {
//		log.Printf("Error sending message: %v", err)
//	}
//	log.Printf("key %v (hash of %v) found in node %v\n", key, hash.Hash(key), resp.IP)
//	upload(resp.IP, key, fmt.Sprintf("%.23s | %s\ndistributed systems is the best :)\n", time.Now().UTC(), key))
//}

//func upload(nodeAddr, key, value string) {
//	//log.Printf("Attempting to lookup string %v", key)
//	var conn *grpc.ClientConn
//	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)
//
//	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("Cannot connect to: %s", err)
//
//	}
//	defer conn.Close()
//
//	req := &gossip.UploadRequest{Key: key, Value: value}
//	client := gossip.NewInternalListenerClient(conn)
//	_, err = client.Upload(context.Background(), req)
//	if err != nil {
//		log.Printf("Error sending message: %v", err)
//	}
//	//log.Printf("key %v uploaded to node %v\n", key, nodeAddr)
//}

func writeExternalFile(nodeAddr, fileName, containerIP string) {
	//log.Printf("Attempting to lookup string %v", key)
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
		log.Fatalf("-- MIKE external fail %s", err)
	}

	log.Printf("\nSuccess upload info to the following chord node: %+v\n", res)
}

func resolveFile(nodeAddr, fileName string) {
	//log.Printf("Attempting to lookup string %v", key)
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
		log.Fatalf("-- MIKE external fail %s", err)
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
		fmt.Println("Sleeping the night away")
		time.Sleep(5 * time.Second)
		writeExternalFile(attachedNode, i, containerIP)
		// fmt.Println("Sleep again, retrieving back the same file in 5s")
		// time.Sleep(5 * time.Second)
		// resolveFile(attachedNode, i)
	}

}
