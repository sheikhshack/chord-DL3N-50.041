package gossip

import (
	"context"
	"errors"
	"log"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

// methods avaiable to the gossiper via the node package
type node interface {
	FindSuccessor(hashed int) string
	GetPredecessor() string
	NotifyHandler(possiblePredecessor string)
}

type Gossiper struct {
	Node node

	pb.UnimplementedInternalListenerServer
}

// This method is server-side
// TODO: consider that we need to spin up goroutine whenever we serve a new request
func (g *Gossiper) Emit(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("received: %+v\n", in)

	switch in.Command {
	case pb.Command_HEALTHCHECK:
		resBool := g.healthcheckHandler()
		res := &pb.Response{
			Body: &pb.Response_Body{
				Success: resBool,
			},
		}
		log.Printf("sending out: %+v\n", res)
		return res, nil
	default:
		return nil, errors.New("CMD not recognised?? use proto.equal???")
	}
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

//Not Used?
// handler to get (Lookup)
func (g *Gossiper) getHandler(key string) ([]byte, error) {
	return store.Get(key)
}
