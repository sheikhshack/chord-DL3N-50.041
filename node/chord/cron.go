package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"math"
)

func (n *Node) stabilize() {
	x := n.Gossiper.GetPredecessor(n.ID, n.successor)
	if hash.IsInRange(hash.Hash(x), hash.Hash(n.ID), hash.Hash(n.successor)) {
		n.setSuccessor(x)
	}
	n.notify(n.successor)
}

//implemented differently from pseudocode, n thinks it might be the predecessor of id
func (n *Node) notify(id string) {
	n.Gossiper.Notify(n.ID, id)
}

// used as a handler func for gossip.Gossiper.NotifyHandler
func (n *Node) NotifyHandler(possiblePredecessor string) {
	//possiblePredecessor is Request's pred
	if (n.GetPredecessor() == "") ||
		(hash.IsInRange(
			hash.Hash(possiblePredecessor),
			hash.Hash(n.GetPredecessor()),
			hash.Hash(n.ID),
		)) {
		n.SetPredecessor(possiblePredecessor)
	}
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next >= len(n.fingers) {
		n.next = 0
	}
	x := int(math.Pow(2, float64(n.next-1)))
	n.fingers[n.next] = n.FindSuccessor(hash.Hash(n.ID) + x)
}

func (n *Node) checkPredecessor() bool {
	res, err := n.Gossiper.Healthcheck(n.ID, n.predecessor)
	if err != nil {
		return false
	}
	return res
}
