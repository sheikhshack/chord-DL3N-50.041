package grpc

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

type Listener struct {
	node *chord.Node
	// each pair of channels is attached to the respective command
	ports    map[Command]*Ports
	AddrBook map[string]map[Command]*Ports
}

type Ports struct {
	req *chan Request
	res *chan Response
}

// handler to findSuccessor
//func (s *Listener) FindSuccessorHandler(key int) (id string) {
func (s *Listener) FindSuccessorHandler() {
	cmd := FindSuccessorCmd
	for req := range *(s.ports[cmd].req) {
		key := req.Body.FindSuccessor.KeySlot
		resChan := *(s.ports[cmd].res)
		//return s.node.FindSuccessor(key)
		id := s.node.FindSuccessor(key)
		res := Response{
			Command:     cmd,
			RequesterID: req.RequesterID,
			TargetID:    req.TargetID,
			Body: ResponseBody{
				FindSuccessor: &IDBody{
					ID: id,
				},
			},
		}
		resChan <- res
	}
}

// handler to join
//func (s *Listener) JoinHandler(k string) string {
func (s *Listener) JoinHandler() string {
	//k= previous node's id
	return s.node.FindSuccessor(chord.Hash(k))
}

//Not Used?
// handler to get (Lookup)
//func (s *Listener) GetHandler(key string) ([]byte, error) {
func (s *Listener) GetHandler() ([]byte, error) {
	return store.Get(key)
}

// handler to healthcheck (checkPredecessor)
//func (s *Listener) HealthcheckHandler() bool {
func (s *Listener) HealthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

//func (s *Listener) GetPredecessorHandler() string {
func (s *Listener) GetPredecessorHandler() string {
	return s.node.GetPredecessor()
}

// notifyHandler handles notify requests and returns if id is req between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
//func (s *Listener) NotifyHandler(possible_pred string) {
func (s *Listener) NotifyHandler() {
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
