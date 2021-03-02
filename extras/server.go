package main


import (
	"fmt"
	"github.com/sheikhshack/distributed-chaos-50.041/basic"
	"log"
	"net"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Basic Server tutorial!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := basic.Listener{}

	grpcServer := grpc.NewServer()

	basic.RegisterBasicServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}