package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/grpc"
)

type Node struct {
	ID          string // maybe IP address
	fingers     []string
	predecessor string
	successor   string
}

// New creates and returns a new Node
func New(id string) Node {
	panic("not implemented")
}

// grpc
func (n *Node) Lookup(k string) []byte {
	//listen on grpc
	//findsuccessor and returns ip
	// grpc call (Get) to another node to retrieve the value
	// feed value to gRPC (bytes, response)
	grpc.Get("4", k)
	panic("not implemented")
}

// grpc
func (n *Node) findSuccessor(k string) string {
	panic("not implemented")
}

func (n *Node) closestPrecedingNode(k string) string {
	panic("not implemented")
}

func (n *Node) initRing() {
	panic("not implemented")
}

// grpc
func (n *Node) join(id string) {
	panic("not implemented")
}
