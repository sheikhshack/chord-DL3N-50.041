package grpc

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

type node interface {
	FindSuccessor(hashed int) string
	GetPredecessor() string
	NotifyHandler(possiblePredecessor string)
}

type Listener struct {
	node *node
}

// handler to findSuccessor
func (s *Listener) FindSuccessorHandler(key int) (id string) {
	return (*s.node).FindSuccessor(key)
}

// handler to join
func (s *Listener) JoinHandler(fromID string) string {
	//fromID= previous node's id
	return (*s.node).FindSuccessor(hash.Hash(fromID))
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
	return (*s.node).GetPredecessor()
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (s *Listener) NotifyHandler(possiblePredecessor string) {
	(*s.node).NotifyHandler(possiblePredecessor)
}
