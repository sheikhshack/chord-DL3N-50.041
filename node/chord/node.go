package chord

import (
	"log"
	"os"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

const SUCCESSOR_LIST_SIZE = 3
const FINGER_TABLE_SIZE = 16

type Node struct {
	ID            string // maybe IP address
	fingers       []string
	predecessor   string
	next          int
	successorList []string

	Gossiper *gossip.Gossiper
}

// New creates and returns a new Node
func New(id string) *Node {
	// 16 is finger table size
	n := &Node{ID: id, next: 0, fingers: make([]string, FINGER_TABLE_SIZE), successorList: make([]string, SUCCESSOR_LIST_SIZE)}

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
	//if n.successor == n.ID {
	//	return n.ID
	//}
	if hash.IsInRange(hashed, hash.Hash(n.ID), hash.Hash(n.GetSuccessor())+1) {
		return n.GetSuccessor()
	} else {
		nPrime := n.closestPrecedingNode(hashed)

		successor, err := n.Gossiper.FindSuccessor(n.ID, nPrime, hashed)
		if err != nil {

			log.Printf("Error in FindSucessor(). Fixing fingerTables.\n")

			for i := 0; i < FINGER_TABLE_SIZE; i++ {
				n.fixFingers()
			}

			// Assume that FingerTables will eventually be corrected if there is >1 node alive
			if n.successorList[0] != n.ID {
				return n.FindSuccessor(hashed)
			} else {
				// No more alive successor nodes except itself (Same as commented out edge case)
				return n.ID
			}
		}
		return successor
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
		log.Println(n.ID, "successor is", n.GetSuccessor(), ", predecessor is", n.predecessor)
		n.checkPredecessor()
		n.stabilize()
		n.fixFingers()
		time.Sleep(time.Millisecond * 1000)
	}
}

// // Fake test
// func (n *Node) cron() {
// 	time.Sleep(time.Millisecond * 10000)
// 	for i := 0; ; i++ {
// 		log.Println(n.ID, "successor is", n.GetSuccessor(), ", predecessor is", n.predecessor)
// 		n.checkPredecessor()
// 		n.stabilize()
// 		n.fixFingers()
// 		time.Sleep(time.Millisecond * 1000)

// 		if i == 15 && n.ID == "charlie" {

// 			// Create for own store
// 			keyString := "store's key 1"
// 			valueString := "store's value 1"
// 			valueByte := []byte(valueString)
// 			store.New(n.ID, keyString, valueByte)

// 			n.ReplicateToSuccessorList(keyString, valueString)

// 		}

// 		if i > 20 {
// 			keyString := "store's key 1"
// 			content, err := store.Get("charlie", keyString)
// 			if err != nil {
// 				log.Println("[Read Fail] Error in finding key.")
// 			} else {
// 				log.Printf("[Read Successful] Value of content is: %s.\n", content)
// 			}
// 		}
// 	}
// }

//change successor
func (n *Node) SetSuccessor(id string) {
	n.successorList[0] = id
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

// exposed version is essentially the same crap
func (n *Node) WriteFileToNode(nodeAddr, fileName, ip string) {

}

func (n *Node) WriteFile(nodeId, fileName, ip string) error {
	key := fileName
	val := ip
	log.Printf("--- FS: Triggering File Write to Chord Node for key [%v] with content %v \n", key, val)

	fileByte := []byte(val)
	output := store.New(nodeId, key, fileByte)

	return output
}

func (n *Node) ReplicateToSuccessorList(fileName, ip string) {

	// Assumes that successorList nodes repeat after it contains own node
	for i := range n.successorList {
		if n.successorList[i] != n.ID {
			n.replicateToNode(n.successorList[i], fileName, ip)
		} else {
			break
		}
	}
}

func (n *Node) replicateToNode(toID, fileName, ip string) bool {
	status, err := n.Gossiper.ReplicateToNode(n.ID, toID, fileName, ip)

	if err != nil {
		log.Printf("Error in replicating file to Node: %s\n", toID)
		return false
	}

	return status

}
