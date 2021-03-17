package exposed

import (
	"context"
	"fmt"
	exposed "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type node interface {


}

type ExternalService struct {
	Node node
	pb.UnimplementedInternalListenerServer


}

func (e *ExternalService) NewServerAndListen(listenPort int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInternalListenerServer(s, g)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Printf("Listening on port %v\n", listenPort)
	}
	return s
}


func (e *ExternalService) Upload (ctx context.Context, uploadRequest *exposed.UploadRequest) (*exposed.IPResponse, error){
	log.Printf("Upload Method triggered \n")
	key := uploadRequest.Key
	val := uploadRequest.Value
	// weeping carry on from here
	return &exposed.IPResponse{IP: "IP"}, nil
}

func (e *ExternalService) LookupIP (ctx context.Context, lookupRequest *exposed.Request) (*exposed.IPResponse, error){
	log.Printf("Lookup Method \n")
	key := lookupRequest.Key
	// weeping carry on from here
	return  &exposed.IPResponse{IP: "IP"},nil
}

func (e *ExternalService) Download (ctx context.Context, downloadRequest *exposed.Request) (*exposed.Response, error) {
	log.Printf("Download Method triggered \n")
	key := downloadRequest.Key
	// weeping carry on from here
	return &exposed.Response{Value: "SOME STUPID STRING"}, nil
}


