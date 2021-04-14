package chord

import (
	"math"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/log"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

func (n *Node) cron() {
	time.Sleep(time.Millisecond * 10000)
	for {
		n.checkPredecessor()
		n.stabilize()
		n.fixFingers()
		time.Sleep(time.Millisecond * 1000)
	}
}

//TODO: Handle the case when the node is in the successorList as well
func (n *Node) stabilize() {
	if n.GetSuccessor() == n.GetID() {
		return
	}

	x, err := n.Gossiper.GetPredecessor(n.GetID(), n.GetSuccessor())
	if err != nil {
		log.Warn.Printf("Successor %v is down: %+v\n", n.GetSuccessor(), err)
		n.fixSuccessorList()
		return
	}

	// Init temp previous successorList variable
	prevSuccessorList := make([]string, n.replicaCount)
	copy(prevSuccessorList, n.GetSuccessorList())

	if hash.IsInRange(hash.Hash(x), hash.Hash(n.GetID()), hash.Hash(n.GetSuccessor())) {
		// Check if x (supposedly predecessor of successor) is alive
		if n.healthCheck(x) {
			n.SetSuccessor(x)
		} else {
			log.Warn.Printf("Predecessor of %s is down. Notifying successor.\n", n.GetSuccessor())
			n.notify(n.GetSuccessor())
			return
		}
	}

	// Get succ list of new successor
	successor_SuccessorList, err := n.Gossiper.GetSuccessorList(n.GetID(), n.GetSuccessor())
	if err != nil {
		log.Warn.Printf("Successor %v is down: %+v\n", n.GetSuccessor(), err)
		n.fixSuccessorList()
		return
	}

	n.updateSuccessorList(successor_SuccessorList, prevSuccessorList)

	n.notify(n.GetSuccessor())
}

func (n *Node) healthCheck(id string) bool {
	if id == "" || id == n.GetID() {
		return true
	}
	res, err := n.Gossiper.Healthcheck(n.GetID(), id)
	if err != nil {
		return false
	}
	return res
}

// TODO: Mutex locks for this
func (n *Node) fixSuccessorList() {
	n.successorList = n.successorList[1:]
	n.successorList = append(n.successorList, "")

	isNewSuccessorAlive := n.healthCheck(n.GetSuccessor())
	if !isNewSuccessorAlive {
		log.Warn.Printf("Successor %v is down.\n", n.GetSuccessor())
		n.fixSuccessorList()

	} else {
		if n.GetSuccessor() != "" {
			n.migrationFault(n.GetSuccessor())
		} else {
			log.Error.Fatalf("Break in chord logical structure. Shutting down node.\n")
			return
		}

	}

}

func (n *Node) updateSuccessorList(successor_SuccessorList []string, prevSuccessorList []string) {
	copy(n.successorList[1:], successor_SuccessorList[:n.replicaCount-1])

	newElements, missingElements := compareList(prevSuccessorList, n.GetSuccessorList())

	if len(newElements) > 0 || len(missingElements) > 0 {
		log.Info.Printf("Updated SuccessorList: %+v\n", n.GetSuccessorList())
		keys, values := stringifyAllLocalFiles()

		if len(newElements) > 0 && keys != "" {
			// Replicate local files to new replica(s)
			n.replicateToNodeList(newElements, keys, values)
		}

		if len(missingElements) > 0 && keys != "" {
			// Delete local files from old replica(s)
			n.deleteFromNodeList(missingElements, keys)
		}
	}

}

//implemented differently from pseudocode, n thinks it might be the predecessor of id
func (n *Node) notify(id string) {
	if id == n.GetID() {
		return
	}
	n.Gossiper.Notify(n.GetID(), id)
}

// used as a handler func for gossip.Gossiper.NotifyHandler
func (n *Node) NotifyHandler(possiblePredecessor string) {
	//possiblePredecessor is Request's pred
	if (n.GetPredecessor() == "") ||
		(hash.IsInRange(
			hash.Hash(possiblePredecessor),
			hash.Hash(n.GetPredecessor()),
			hash.Hash(n.GetID()),
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
	n.fingers[n.next] = n.FindSuccessor(hash.Hash(n.GetID()) + x)
}

func (n *Node) checkPredecessor() {

	if !n.healthCheck(n.predecessor) {
		log.Warn.Printf("Predecessor %s is down.\n", n.predecessor)
		n.SetPredecessor("")
	}
}
