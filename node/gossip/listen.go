package gossip

import (
	"context"
	"errors"
	"log"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

type node interface {
	FindSuccessor(hashed int) string
	GetPredecessor() string
	NotifyHandler(possiblePredecessor string)
}

type Listener struct {
	Node node
	pb.UnimplementedInternalListenerServer
}

func (s *Listener) Emit(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("received: %+v\n", in)

	switch in.Command {
	case pb.Command_HEALTHCHECK:
		resBool := s.HealthcheckHandler()
		res := &pb.Response{
			Body: &pb.Response_Body{
				Healthcheck: &pb.Response_SuccessBody{Success: resBool},
			},
		}
		log.Printf("sending out: %+v\n", res)
		return res, nil
	default:
		return nil, errors.New("CMD not recognised?? use proto.equal???")
	}
}

// handler to findSuccessor
func (s *Listener) FindSuccessorHandler(key int) (id string) {
	return s.Node.FindSuccessor(key)
}

// handler to join
func (s *Listener) JoinHandler(fromID string) string {
	//fromID= previous node's id
	return s.Node.FindSuccessor(hash.Hash(fromID))
}

//Not Used?
// handler to get (Lookup)
func (s *Listener) GetHandler(key string) ([]byte, error) {
	return store.Get(key)
}

// handler to healthcheck (checkPredecessor)
func (s *Listener) HealthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

func (s *Listener) GetPredecessorHandler() string {
	return s.Node.GetPredecessor()
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (s *Listener) NotifyHandler(possiblePredecessor string) {
	s.Node.NotifyHandler(possiblePredecessor)
}
