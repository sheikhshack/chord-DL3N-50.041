package chord

import "github.com/sheikhshack/distributed-chaos-50.041/node/store"

/* All these functions are handler functions to the grpc
 * (perspective: node is listening)
 * can replace with/move to client dir or smth when implementing grpc proper
 */

// handler to findSuccessor
func (n *Node) findSuccessorHandler(key int) (id string) {
	return n.findSuccessor(key)
}

// handler to join
func (n *Node) joinHandler(k string) string {
	//k= previous node's id
	return n.findSuccessor(Hash(k))
}

//Not Used?
// handler to get (Lookup)
func (n *Node) getHandler(key string) ([]byte, error) {
	return store.Get(key)
}

// handler to healthcheck (checkPredecessor)
func (n *Node) healthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

func (n *Node) getPredecessorHandler() string {
	return n.predecessor
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (n *Node) notifyHandler(possible_pred string) {
	//possible_pred is Request's pred
	if (n.predecessor == "") || (IsInRange(Hash(possible_pred), Hash(n.predecessor), Hash(n.ID))) {
		n.setPredecessor(possible_pred)
	}
}
