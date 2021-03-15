package chord

/* All these functions are functions to handle grpc calls
 * (perspective: node is sending grpc requests)
 * can replace with/move to client dir or smth when implementing grpc proper
 * this and listen.go can be inside another package actually (grpc package?)
 */

// called by findSuccessor
func (n *Node) findSuccessorRequest(id string, key int) string {
	panic("not implemented")
}

// called by join
func (n *Node) joinRequest(id string) string {
	//k = n.ID
	panic("not implemented")
}

// Not used?
// Called by Lookup
func (n *Node) getRequest(key string) ([]byte, error) {
	panic("not implemented")
}

// called by checkPredecessor
func (n *Node) healthcheckRequest() bool {
	panic("not implemented")
}

//Get the predecessor of the node
func (n *Node) getPredecessorRequest(id string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the predecessor of id
func (n *Node) notifyRequest(id string) bool {
	//pred = n.ID
	panic("not implemented")
}
