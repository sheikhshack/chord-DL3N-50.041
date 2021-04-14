package chord

//change successor
func (n *Node) SetSuccessor(id string) {
	n.successorList[0] = id
}

//change predecessor
func (n *Node) SetPredecessor(id string) {
	n.predecessor = id
}

// get predecessor
func (n *Node) GetPredecessor() string {
	return n.predecessor
}

// Get Successor List
func (n *Node) GetSuccessorList() []string {
	return n.successorList
}

func (n *Node) GetSuccessor() string {
	return n.successorList[0]
}

func (n *Node) GetID() string {
	return n.ID
}

func (n *Node) GetFingers() []string {
	return n.fingers
}
