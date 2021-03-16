package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/grpc"
	"math"
)

func (n *Node) stabilize() {
	x := grpc.GetPredecessor(n.ID, n.successor)
	if IsInRange(Hash(x), Hash(n.ID), Hash(n.successor)) {
		n.setSuccessor(x)
	}
	n.notify(n.successor)
}

//implemented differently from pseudocode, n thinks it might be the predecessor of id
func (n *Node) notify(id string) {
	grpc.Notify(n.ID, id)
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next >= len(n.fingers) {
		n.next = 0
	}
	x := int(math.Pow(2, float64(n.next-1)))
	n.fingers[n.next] = n.FindSuccessor(Hash(n.ID) + x)
}

func (n *Node) checkPredecessor() bool {
	return grpc.Healthcheck(n.ID, n.predecessor)
}
