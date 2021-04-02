package chord

import (
	"log"
	"math"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

//TODO: [fault] When shutting down alpha node, will crash its predecessor occasionally
//TODO: Handle the case when the node is in the successorList as well
func (n *Node) stabilize() {
	log.Println("Stabilizing", n.ID)
	if n.successorList[0] == n.ID {
		return
	}

	// Set Predecessor to "" when predecessor is down
	if !n.checkPredecessor() {
		log.Printf("%s's Predecessor is down.\n", n.ID)
		n.SetPredecessor("")
	}

	x, err := n.Gossiper.GetPredecessor(n.ID, n.successorList[0])
	// log.Println("node ID: ", n.ID)
	// log.Println("node successorList[0]: ", n.successorList[0])
	// log.Println("Value of x: ", x)

	if err != nil {
		log.Printf("error in stabilize[GetPredecessor]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	// Check if x (supposedly predecessor of successor) is alive (decouple)
	if n.healthCheck(x) {
		if hash.IsInRange(hash.Hash(x), hash.Hash(n.ID), hash.Hash(n.successorList[0])) {
			n.SetSuccessor(x)
		}
	} else {
		log.Printf("Predecessor of %s is down. Notify the successor.\n", n.successorList[0])
		n.notify(n.successorList[0])
		return
	}

	// Get succ list of new successor
	succSuccList, err := n.Gossiper.GetSuccessorList(n.ID, n.successorList[0])

	if err != nil {
		log.Printf("error in stabilize[GetSuccessorList]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	n.updateSuccessorList(n.successorList, succSuccList)

	log.Println("Value of successorList: ", n.successorList)

	n.notify(n.successorList[0])
}

func (n *Node) healthCheck(id string) bool {
	res, err := n.Gossiper.Healthcheck(n.ID, id)
	if err != nil {
		return false
	}
	return res
}

func (n *Node) fixSuccessorList() {

	if len(n.successorList) <= 1 {
		n.successorList = make([]string, SUCCESSOR_LIST_SIZE)
	} else {
		n.successorList = n.successorList[1:]
	}

}

func (n *Node) updateSuccessorList(succList []string, succSuccList []string) {

	tempSuccList := succList[:1]
	tempSuccList = append(tempSuccList, succSuccList[:len(succSuccList)-1]...)

	n.successorList = tempSuccList
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

func (n *Node) checkPredecessor() bool {
	return n.healthCheck(n.predecessor)
}
