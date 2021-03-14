package chord

/* All these functions are functions to handle grpc calls
 * (perspective: node is sending grpc requests)
 * can replace with/move to client dir or smth when implementing grpc proper
 * this and listen.go can be inside another package actually (grpc package?)
 */

// called by findSuccessor
func (n *Node) findSuccessorRequest(key string) (id string) {
	panic("not implemented")
}

// called by join
func (n *Node) joinRequest(id string) {
	panic("not implemented")
}

// called by Lookup
func (n *Node) getRequest(key string) ([]byte, error) {
	panic("not implemented")
}

// called by checkPredecessor
func (n *Node) healthcheckRequest() bool {
	panic("not implemented")
}

// called by notify
func (n *Node) notifyRequest(id string) bool {
	panic("not implemented")
}
