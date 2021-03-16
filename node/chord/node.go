package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

type Node struct {
	ID          string // maybe IP address
	fingers     []string
	predecessor string
	successor   string
	next        int

	Gossiper *gossip.Gossiper
}

// New creates and returns a new Node
func New(id string) *Node {
	n := &Node{ID: id}
	n.Gossiper = &gossip.Gossiper{
		Node: n,
	}
	return n
}

func (n *Node) Lookup(k string) (ip string) {
	//listen on gossip
	//findsuccessor and returns ip
	return n.FindSuccessor(hash.Hash(k))
}

func (n *Node) FindSuccessor(hashed int) string {
	if hash.IsInRange(hashed, hash.Hash(n.ID), hash.Hash(n.successor)+1) {
		return n.successor
	} else {
		n_prime := n.closestPrecedingNode(hashed)
		return n.Gossiper.FindSuccessor(n.ID, n_prime, hashed)
	}
}

//searches local table for highest predecessor of id
func (n *Node) closestPrecedingNode(hashed int) string {
	m := len(n.fingers)
	for i := m; i > 0; i-- {
		if hash.IsInRange(hash.Hash(n.fingers[i]), hash.Hash(n.ID), hashed) {
			return n.fingers[i]
		}
	}
	return n.ID
}

func (n *Node) InitRing() {
	n.SetPredecessor("")
	n.setSuccessor(n.ID)
}

func (n *Node) join(id string) {
	successor := n.Gossiper.Join(n.ID, id)
	n.SetPredecessor("")
	n.setSuccessor(successor)
}

//change successor
func (n *Node) setSuccessor(id string) {
	n.successor = id
}

//change predecessor
func (n *Node) SetPredecessor(id string) {
	//TODO: Need to have the case where id is ""
	n.predecessor = id
}

// get predecessor
func (n *Node) GetPredecessor() string {
	return n.predecessor
}
