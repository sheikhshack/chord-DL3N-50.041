package grpc

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

type Listener struct {
	node *chord.Node
}

// handler to findSuccessor
func (s *Listener) FindSuccessorHandler(key int) (id string) {
	return s.node.FindSuccessor(key)
}

// handler to join
func (s *Listener) JoinHandler(k string) string {
	//k= previous node's id
	return s.node.FindSuccessor(chord.Hash(k))
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
	return s.node.GetPredecessor()
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (s *Listener) NotifyHandler(possible_pred string) {
	//possible_pred is Request's pred
	if (s.node.GetPredecessor() == "") ||
		(chord.IsInRange(
			chord.Hash(possible_pred),
			chord.Hash(s.node.GetPredecessor()),
			chord.Hash(s.node.ID),
		)) {
		s.node.SetPredecessor(possible_pred)
	}
}
