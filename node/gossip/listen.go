package gossip

import (
	"context"
	"errors"
	"fmt"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

// methods avaiable to the gossiper via the node package
type node interface {
	FindSuccessor(hashed int) string
	GetPredecessor() string
	NotifyHandler(possiblePredecessor string)
	LookupIP(k string) (ip string)
	GetID() (id string)
}

type Gossiper struct {
	Node node

	pb.UnimplementedInternalListenerServer
}

func (g *Gossiper) NewServerAndListen(listenPort int) *grpc.Server {
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

// This method is server-side
// TODO: consider that we need to spin up goroutine whenever we serve a new request
func (g *Gossiper) Emit(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("received: %+v\n", in)
	var res *pb.Response

	switch in.Command {
	case pb.Command_FIND_SUCCESSOR:
		key := int(in.GetBody().HashSlot)
		id := g.findSuccessorHandler(key)
		res = &pb.Response{
			Command:     pb.Command_FIND_SUCCESSOR,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: id,
			},
		}

	case pb.Command_JOIN:
		fromID := in.GetRequesterID()
		id := g.joinHandler(fromID)
		res = &pb.Response{
			Command:     pb.Command_JOIN,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: id,
			},
		}

	case pb.Command_HEALTHCHECK:
		isHealthy := g.healthcheckHandler()
		res = &pb.Response{
			Command:     pb.Command_HEALTHCHECK,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				IsHealthy: isHealthy,
			},
		}

	case pb.Command_GET_PREDECESSOR:
		predecessorID := g.getPredecessorHandler()
		res = &pb.Response{
			Command:     pb.Command_GET_PREDECESSOR,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: predecessorID,
			},
		}

	case pb.Command_NOTIFY:
		possiblePredecessor := in.GetRequesterID()
		g.notifyHandler(possiblePredecessor)
		res = &pb.Response{
			Command:     pb.Command_NOTIFY,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body:        &pb.Response_Body{},
		}

	default:
		return nil, errors.New("command not recognised")
	}

	log.Printf("sending out: %+v\n", res)
	return res, nil
}

// handler to findSuccessor
func (g *Gossiper) findSuccessorHandler(key int) (id string) {
	return g.Node.FindSuccessor(key)
}

// handler to join
func (g *Gossiper) joinHandler(fromID string) string {
	//fromID= previous node's id
	return g.Node.FindSuccessor(hash.Hash(fromID))
}

// handler to healthcheck (checkPredecessor)
func (g *Gossiper) healthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

func (g *Gossiper) getPredecessorHandler() string {
	return g.Node.GetPredecessor()
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (g *Gossiper) notifyHandler(possiblePredecessor string) {
	g.Node.NotifyHandler(possiblePredecessor)
}

func (g *Gossiper) Upload(ctx context.Context, uploadRequest *pb.UploadRequest) (*pb.UploadResponse, error) {
	log.Printf("Upload Method triggered \n")
	key := uploadRequest.Key
	val := uploadRequest.Value

	fileByte := []byte(val)
	output := store.New(key, fileByte)

	return &pb.UploadResponse{IP: g.Node.GetID()}, output
}

func (g *Gossiper) CheckIP(ctx context.Context, lookupRequest *pb.CheckRequest) (*pb.CheckResponse, error) {
	log.Printf("Lookup Method \n")
	key := lookupRequest.Key
	ip := g.Node.LookupIP(key)
	return &pb.CheckResponse{IP: ip}, nil
}

func (g *Gossiper) Download(ctx context.Context, downloadRequest *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	log.Printf("Download Method triggered \n")
	key := downloadRequest.Key

	fileByte, status := store.Get(key)
	if status == nil {
		v := string(fileByte)
		return &pb.DownloadResponse{Value: v}, nil
	} else {
		return nil, status
	}
}