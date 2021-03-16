package gossip

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
)

// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

const (
	LISTEN_PORT = 9000
	EMIT_PORT   = 9001
)

func (g *Gossiper) emit(nodeAddr string, request *pb.Request) (*pb.Response, error) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)
	//log.Printf("%v\n", )
	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.Emit(context.Background(), request)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return nil, err
	}
	log.Printf("Response from server: %s", response.Body)

	return response, nil
}

// called by FindSuccessor
func (g *Gossiper) FindSuccessor(fromID, toID string, key int) string {
	panic("not implemented")
}

// called by join
func (g *Gossiper) Join(fromID, toID string) string {
	//k = n.ID
	panic("not implemented")
}

// called by checkPredecessor
// TODO: change all method signatures to include error in return
func (g *Gossiper) Healthcheck(fromID, toID string) (bool, error) {
	req := &pb.Request{
		Command:     pb.Command_HEALTHCHECK,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{Healthcheck: &pb.Request_NullBody{}},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return false, err
	}
	return res.GetBody().GetHealthcheck().GetSuccess(), nil
}

//Get the predecessor of the node
func (g *Gossiper) GetPredecessor(fromID, toID string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the predecessor of id
func (g *Gossiper) Notify(fromID, toID string) {
	//pred = n.ID
	panic("not implemented")
}

// Not used?
// Called by Lookup
// TODO: move this method to exposed API package
func (g *Gossiper) Get(fromID, toID, key string) ([]byte, error) {
	panic("not implemented")
}
