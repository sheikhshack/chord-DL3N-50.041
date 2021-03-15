package chord

import (
	"math"
)

func (n *Node) stabilize() {
	x := n.getPredecessorRequest(n.successor)
	if IsInRange(Hash(x), Hash(n.ID), Hash(n.successor)) {
		n.setSuccessor(x)
	}
	n.notify(n.successor)
}

// grpc
//implemented differently from pseudocode, n thinks it might be the Predecessor of id
func (n *Node) notify(id string) {
	n.notifyRequest(id)
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next >= len(n.fingers) {
		n.next = 0
	}
	x := int(math.Pow(2, float64(n.next-1)))
	n.fingers[n.next] = n.FindSuccessor(Hash(n.ID) + x)
}

// grpc (healthcheck)
func (n *Node) checkPredecessor() bool {
	return n.healthcheckRequest()
}
