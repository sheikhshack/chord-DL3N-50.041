package chord

import (
	"log"
	"os"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

type Node struct {
	ID          string // maybe IP address
	fingers     []string
	predecessor string
	successor   string
	next        int

	Gossiper *gossip.Gossiper
}

// New creates and returns a new Node
func New(id string) *Node {
	n := &Node{ID: id, next: 0}

	if os.Getenv("DEBUG") == "debug" {
		n.Gossiper = &gossip.Gossiper{
			Node:      n,
			DebugMode: true,
		}
	} else {
		n.Gossiper = &gossip.Gossiper{
			Node:      n,
			DebugMode: false,
		}
	}

	return n
}

func (n *Node) InitRing() {
	n.SetPredecessor("")
	n.SetSuccessor(n.ID)
	go n.cron()
}

func (n *Node) Join(id string) {
	successor, err := n.Gossiper.Join(n.ID, id)
	if err != nil {
		// TODO: handle this error
		// we can pass the error back and have main.go to exit gracefully with helpful message
		log.Fatalf("error in join: %+v\n", err)
	}
	n.SetPredecessor("")
	n.SetSuccessor(successor)
	go n.cron()
}

func (n *Node) FindSuccessor(hashed int) string {
	// edge case of having only one node in ring
	if n.successor == n.ID {
		return n.ID
	}
	if hash.IsInRange(hashed, hash.Hash(n.ID), hash.Hash(n.successor)+1) {
		return n.successor
	} else {
		//n_prime := n.closestPrecedingNode(hashed)
		//if n_prime == n.ID {
		//	return n.ID
		//}
		//successor, err := n.Gossiper.FindSuccessor(n.ID, n_prime, hashed)
		successor, err := n.Gossiper.FindSuccessor(n.ID, n.successor, hashed)
		if err != nil {
			// TODO: handle this error
			log.Fatalf("error in FindSuccessor: %+v\n", err)
		}
		return successor
	}
}

//searches local table for highest predecessor of id
func (n *Node) closestPrecedingNode(hashed int) string {
	m := len(n.fingers)
	for i := m; i > 0; i-- {
		if hash.IsInRange(hash.Hash(n.fingers[i]), hash.Hash(n.ID), hashed) {
			return n.fingers[i]
		}
	}
	return n.ID
}

func (n *Node) cron() {
	time.Sleep(time.Millisecond * 10000)
	for {
		log.Println(n.ID, "successor is", n.successor, ", predecessor is", n.predecessor)
		n.stabilize()
		time.Sleep(time.Millisecond * 1000)

		// if n.ID == "alpha" {
		// 	log.Println("LOOK HEREEEEEEEEEEEEEEEE -> Lookup('hello')'s sucessor is ", n.Lookup("hello"))
		// }
	}
}

//change successor
func (n *Node) SetSuccessor(id string) {
	n.successor = id
}

//change predecessor
func (n *Node) SetPredecessor(id string) {
	//TODO: Need to have the case where id is ""
	n.predecessor = id
}

// get predecessor
func (n *Node) GetPredecessor() string {
	return n.predecessor
}

func (n *Node) GetSuccessor() string {
	return n.successor
}

func (n *Node) GetID() string {
	return n.ID
}
