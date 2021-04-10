package chord

import (
	"log"
	"math"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

//TODO: Handle the case when the node is in the successorList as well
func (n *Node) stabilize() {
	log.Println("Stabilizing", n.ID)
	if n.GetSuccessor() == n.ID {
		return
	}

	x, err := n.Gossiper.GetPredecessor(n.ID, n.GetSuccessor())
	// log.Println("node ID: ", n.ID)
	// log.Println("node successorList[0]: ", n.GetSuccessor())
	// log.Println("Value of x: ", x)

	if err != nil {
		log.Printf("error in stabilize[GetPredecessor]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	if hash.IsInRange(hash.Hash(x), hash.Hash(n.ID), hash.Hash(n.GetSuccessor())) {
		// Check if x (supposedly predecessor of successor) is alive
		if n.healthCheck(x) {
			n.SetSuccessor(x)
		} else {
			log.Printf("Predecessor of %s is down. Notify the successor.\n", n.GetSuccessor())
			n.notify(n.GetSuccessor())
			return
		}
	}

	// Get succ list of new successor
	succSuccList, err := n.Gossiper.GetSuccessorList(n.ID, n.GetSuccessor())

	if err != nil {
		log.Printf("error in stabilize[GetSuccessorList]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	n.updateSuccessorList(succSuccList)

	log.Println("Value of successorList: ", n.successorList)

	n.notify(n.GetSuccessor())
}

func (n *Node) healthCheck(id string) bool {
	res, err := n.Gossiper.Healthcheck(n.ID, id)
	if err != nil {
		return false
	}
	return res
}

// TODO: Mutex locks for this
func (n *Node) fixSuccessorList() {
	if len(n.successorList) <= 1 {
		n.successorList = make([]string, SUCCESSOR_LIST_SIZE)
	} else {
		n.successorList = n.successorList[1:]
		n.successorList = append(n.successorList, "")
	}
}

// TODO: SuccessorList will have duplicates (Might want to take note for replication) Might have to do away with copy with we want different size
func (n *Node) updateSuccessorList(succSuccList []string) {
	copy(n.successorList[1:], succSuccList[:SUCCESSOR_LIST_SIZE-1])
}

//implemented differently from pseudocode, n thinks it might be the predecessor of id
func (n *Node) notify(id string) {
	if id == n.ID {
		return
	}
	n.Gossiper.Notify(n.ID, id)
}

// used as a handler func for gossip.Gossiper.NotifyHandler
func (n *Node) NotifyHandler(possiblePredecessor string) {
	// log.Printf("[NOTIFY HANDLER] Possible Predecessor of %s is %s.\n", n.ID, possiblePredecessor)
	//possiblePredecessor is Request's pred
	if (n.GetPredecessor() == "") ||
		(hash.IsInRange(
			hash.Hash(possiblePredecessor),
			hash.Hash(n.GetPredecessor()),
			hash.Hash(n.ID),
		)) {
		// log.Printf("[NOTIFY HANDLER - Set Predecessor] Predecessor of %s is %s.\n", n.ID, possiblePredecessor)
		n.SetPredecessor(possiblePredecessor)
	}
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next >= cap(n.fingers) {
		n.next = 0
	}
	x := int(math.Pow(2, float64(n.next)))
	n.fingers[n.next] = n.FindSuccessor(hash.Hash(n.ID) + x)
}

func (n *Node) checkPredecessor() {

	if !n.healthCheck(n.predecessor) {
		log.Printf("%s's Predecessor is down.\n", n.ID)

		// n.migratePredecessorFiles(n.GetPredecessor())
		n.SetPredecessor("")

	}
}

// // Transfer all the files from predecessor to its own folder
// func (n *Node) migratePredecessorFiles(predecessorID string) {

// 	files, err := store.GetAll(predecessorID)

// 	if err != nil {
// 		log.Printf("Error in obtaining all files info in store: %+v\n", err)
// 		return
// 	}

// 	for _, file := range files {
// 		err := store.Migrate(predecessorID, n.ID, file.Name())
// 		if err != nil {
// 			log.Printf("Error in migrating file %s: %+v\n", file.Name(), err)
// 		}
// 	}
// }
