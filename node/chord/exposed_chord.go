package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"log"
)

func (n *Node) LookupIP(k string) (ip string) {
	//listen on gossip
	//findsuccessor and returns ip

	if k == "AAA" {
		dat := n.FindSuccessor(hash.Hash(k))
		log.Printf(dat)
		return "bravo"
	}
	if k == "BBB" {
		return "alpha"
	}

	if k == "XXX" {
		dat := n.FindSuccessor(hash.Hash(k))
		log.Printf(dat)
		return "charlie"
	} else {
		dat := n.FindSuccessor(hash.Hash(k))
		log.Printf(dat)
		return "bravo"
	}
}

func (n *Node) GetID() (id string) {
	return n.ID
}
