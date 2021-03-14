package chord

import "github.com/sheikhshack/distributed-chaos-50.041/node/store"

/* All these functions are handler functions to the grpc
 * (perspective: node is listening)
 * can replace with/move to client dir or smth when implementing grpc proper
 */

// handler to findSuccessor
func (n *Node) findSuccessorHandler(key string) (id string) {
	return n.findSuccessor(key)
}

// handler to join
func (n *Node) joinHandler(id string) {
	panic("not implemented")
}

// handler to get (Lookup)
func (n *Node) getHandler(key string) ([]byte, error) {
	return store.Get(key)
}

// handler to healthcheck (checkPredecessor)
func (n *Node) healthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (n *Node) notifyHandler(id string) bool {
	panic("not implemented")
}
