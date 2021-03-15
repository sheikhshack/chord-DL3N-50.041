package chord

import "github.com/sheikhshack/distributed-chaos-50.041/node/grpc"

type Node struct {
	ID          string // maybe IP address
	fingers     []string
	predecessor string
	successor   string
	next        int

	Listener *grpc.Listener
}

// New creates and returns a new Node
func New(id string) Node {
	n := Node{ID: id}
	return n
}

// grpc
func (n *Node) Lookup(k string) (ip string) {
	//listen on grpc
	//findsuccessor and returns ip
	return n.FindSuccessor(Hash(k))

}

// grpc
func (n *Node) FindSuccessor(hashed int) string {
	if IsInRange(hashed, Hash(n.ID), Hash(n.successor)+1) {
		return n.successor
	} else {
		n_prime := n.closestPrecedingNode(hashed)
		return grpc.FindSuccessor(n_prime, hashed)
	}
}

//searches local table for highest predecessor of id
func (n *Node) closestPrecedingNode(hashed int) string {
	m := len(n.fingers)
	for i := m; i > 0; i-- {
		if IsInRange(Hash(n.fingers[i]), Hash(n.ID), hashed) {
			return n.fingers[i]
		}
	}
	return n.ID
}

func (n *Node) initRing() {
	n.SetPredecessor("")
	n.setSuccessor(n.ID)
}

// grpc
func (n *Node) join(id string) {
	successor := grpc.Join(n.ID, id)
	n.SetPredecessor("")
	n.setSuccessor(successor)
}

//change successor
func (n *Node) setSuccessor(id string) {
	panic("not implemented")
}

//change predecessor
func (n *Node) SetPredecessor(id string) {
	//TODO: Need to have the case where id is ""
	panic("not implemented")
}

// get predecessor
func (n *Node) GetPredecessor() string {
	return n.predecessor
}
