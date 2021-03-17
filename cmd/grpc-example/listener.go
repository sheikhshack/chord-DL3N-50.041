package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
)

func main() {
	log.Printf("Starting Listener!!\n")
	server := chord.New("localhost")
	server.InitRing()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gossip.LISTEN_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInternalListenerServer(s, server.Gossiper)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
