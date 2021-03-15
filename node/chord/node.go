package chord

type Node struct {
	ID          string // maybe IP address
	fingers     []string
	predecessor string
	successor   string
	next        int
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
	return n.findSuccessor(Hash(k))

}

// grpc
func (n *Node) findSuccessor(hashed int) string {
	if IsInRange(hashed, Hash(n.ID), Hash(n.successor)+1) {
		return n.successor
	} else {
		n_prime := n.closestPrecedingNode(hashed)
		return n.findSuccessorRequest(n_prime, hashed)
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
	n.setPredecessor("")
	n.setSuccessor(n.ID)
}

// grpc
func (n *Node) join(id string) {
	successor := n.joinRequest(id)
	n.setPredecessor("")
	n.setSuccessor(successor)
}

//change successor
func (n *Node) setSuccessor(id string) {
	panic("not implemented")
}

//change predecessor
func (n *Node) setPredecessor(id string) {
	//TODO: Need to have the case where id is ""
	panic("not implemented")
}
