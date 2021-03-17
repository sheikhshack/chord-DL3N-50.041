package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

func (n *Node) LookupIP(k string) (ip string) {
	//listen on gossip
	//findsuccessor and returns ip
	return n.FindSuccessor(hash.Hash(k))
}

func (n *Node) GetID() (id string) {
	return n.ID
}
