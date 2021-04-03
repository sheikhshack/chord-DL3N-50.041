package chord

import (
	"log"
	"os"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
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
	// 16 is finger table size
	n := &Node{ID: id, next: 0, fingers: make([]string, 16)}
	log.Printf("Node Hash id: %v\n", hash.Hash(n.ID))

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
	files, err := store.GetAll()
	if err != nil {
		print(err)
	}

	for _, i := range files {
		log.Printf("Filename:%v, HashedFile: %v", i.Name(), hash.Hash(i.Name()))
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
	n.migrationInit(n.successor)
	//edge case of in the 1s window, the node's ideal pred hasn't recognised this node
	go n.cron()

}

func (n *Node) FindSuccessor(hashed int) string {
	// edge case of having only one node in ring
	//if n.successor == n.ID {
	//	return n.ID
	//}
	if hash.IsInRange(hashed, hash.Hash(n.ID), hash.Hash(n.successor)+1) {
		return n.successor
	} else {
		nPrime := n.closestPrecedingNode(hashed)
		successor, err := n.Gossiper.FindSuccessor(n.ID, nPrime, hashed)
		if err != nil {
			// TODO: handle this error
			log.Fatalf("error in FindSuccessor: %+v\n", err)
		}
		return successor
	}
}

//ask successor to migrate files that belong to current node
func (n *Node) migrationInit(successor string) {
	n.Gossiper.MigrationRequestFromNode(successor)
}

//predecessor has asked to migrate files from current nodes
func (n *Node) MigrationHandler(pred string) {

	// Get all the files in the store
	files, err := store.GetAll()
	if err != nil {
		print(err)
	}

	//Loop through them and write over the ones that do not lie in between pred and current node, and then delete if the write is successful
	for _, i := range files {
		log.Printf("Filename:%v, HashedFile: %v", i.Name(), hash.Hash(i.Name()))
		if !hash.IsInRange(hash.Hash(i.Name()), hash.Hash(pred), hash.Hash(n.ID)) {
			val, _ := store.Get(i.Name())
			_, err := n.Gossiper.WriteFileToNode(pred, i.Name(), string(val))
			if err != nil {
				log.Printf("Error in writing file")
			} else {
				store.Delete(i.Name())
			}
		}

	}

}

//searches local table for highest predecessor of id
func (n *Node) closestPrecedingNode(hashed int) string {
	m := cap(n.fingers) - 1
	for i := m; i > 0; i-- {
		if n.fingers[i] == "" {
			continue
		}
		if hash.IsInRange(hash.Hash(n.fingers[i]), hash.Hash(n.ID), hashed) {
			return n.fingers[i]
		}
	}
	return n.ID
}

func (n *Node) cron() {
	time.Sleep(time.Millisecond * 10000)
	for {
		// log.Println(n.ID, "successor is", n.successor, ", predecessor is", n.predecessor)
		n.stabilize()
		n.fixFingers()
		time.Sleep(time.Millisecond * 1000)
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

func (n *Node) GetFingers() []string {
	return n.fingers
}
