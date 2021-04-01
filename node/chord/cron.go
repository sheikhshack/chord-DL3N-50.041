package chord

import (
	"log"
	"math"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

func (n *Node) stabilize() {
	log.Println("Stabilizing", n.ID)
	if n.successorList[0] == n.ID {
		return
	}
	x, err := n.Gossiper.GetPredecessor(n.ID, n.successorList[0])
	// log.Println("node ID: ", n.ID)
	// log.Println("node successorList[0]: ", n.successorList[0])
	// log.Println("Value of x: ", x)

	// TODO: Fix connection between node and new successor node when previous successor node is down
	if err != nil {
		log.Printf("error in stabilize[GetPredecessor]: %+v\n", err)

		// TODO: External function
		if len(n.successorList) <= 1 {
			n.successorList = make([]string, SUCCESSOR_LIST_SIZE)
		} else {
			n.successorList = n.successorList[1:]
		}
	}

	if hash.IsInRange(hash.Hash(x), hash.Hash(n.ID), hash.Hash(n.successorList[0])) {
		n.SetSuccessor(x)
	}

	// Get succ list of new successor
	succSuccList, err := n.Gossiper.GetSuccessorList(n.ID, n.successorList[0])

	if err != nil {
		log.Printf("error in stabilize[GetSuccessorList]: %+v\n", err)
		return
	}

	n.updateSuccessorList(n.successorList, succSuccList)

	log.Println("Value of successorList: ", n.successorList)

	n.notify(n.successorList[0])
}

// First 2 nodes joined => SuccessorList is not accurate (r < n-1 and values of r are different)
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
	//possiblePredecessor is Request's pred
	if (n.GetPredecessor() == "") ||
		(hash.IsInRange(
			hash.Hash(possiblePredecessor),
			hash.Hash(n.GetPredecessor()),
			hash.Hash(n.ID),
		)) {
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
	res, err := n.Gossiper.Healthcheck(n.ID, n.predecessor)
	if err != nil {
		return false
	}
	return res
}
