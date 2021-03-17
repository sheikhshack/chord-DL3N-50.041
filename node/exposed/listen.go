package exposed

import (
	"context"
	"fmt"
	exposed "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
	"google.golang.org/grpc"
	"log"
	"net"
)

type node interface {
	LookupIP(k string) (ip string)
	GetID() (id string)
}

type ExternalService struct {
	Node node
	exposed.UnimplementedExternalListenerServer
}

func (e *ExternalService) NewServerAndListen(listenPort int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	exposed.RegisterExternalListenerServer(s, e)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Printf("Listening on port %v\n", listenPort)
	}
	return s
}

func (e *ExternalService) Upload(ctx context.Context, uploadRequest *exposed.UploadRequest) (*exposed.IPResponse, error) {
	log.Printf("Upload Method triggered \n")
	key := uploadRequest.Key
	val := uploadRequest.Value

	fileByte := []byte(val)
	output := store.New(key, fileByte)

	return &exposed.IPResponse{IP: e.Node.GetID()}, output
}

func (e *ExternalService) CheckIP(ctx context.Context, lookupRequest *exposed.Request) (*exposed.IPResponse, error) {
	log.Printf("Lookup Method \n")
	key := lookupRequest.Key
	ip := e.Node.LookupIP(key)
	return &exposed.IPResponse{IP: ip}, nil
}

func (e *ExternalService) Download(ctx context.Context, downloadRequest *exposed.Request) (*exposed.Response, error) {
	log.Printf("Download Method triggered \n")
	key := downloadRequest.Key

	fileByte, status := store.Get(key)
	if status == nil {
		v := string(fileByte)
		return &exposed.Response{Value: v}, nil
	} else {
		return nil, status
	}
}
