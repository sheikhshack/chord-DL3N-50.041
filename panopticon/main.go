package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
)

type Tower struct {
	data *sync.Map

	pb.UnimplementedInternalListenerServer
}

func (t *Tower) NewServerAndListen(listenPort int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInternalListenerServer(s, t)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Printf("Listening on port %v\n", listenPort)
	}
	return s
}

func (t *Tower) Debug(ctx context.Context, in *pb.DebugMessage) (*pb.DebugResponse, error) {
	data := nodeData{
		nodeID:      in.GetFromID(),
		predecessor: in.GetPredecessor(),
		successor:   in.GetSuccessor(),
		fingers:     in.GetFingers(),
	}
	t.data.Store(in.GetFromID(), data)

	return &pb.DebugResponse{Success: true}, nil
}

func main() {
	log.Printf("hello world??")
	data := &sync.Map{}
	panopticon := Tower{data: data}
	go panopticon.display()
	panopticon.NewServerAndListen(9000)
}
