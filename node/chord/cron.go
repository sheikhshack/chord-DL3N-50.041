package chord

import (
	"log"
	"math"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

//TODO: Handle the case when the node is in the successorList as well
func (n *Node) stabilize() {
	log.Println("Stabilizing", n.GetID())
	if n.GetSuccessor() == n.GetID() {
		return
	}

	x, err := n.Gossiper.GetPredecessor(n.GetID(), n.GetSuccessor())

	if err != nil {
		log.Printf("error in stabilize[GetPredecessor]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	// Init temp previous successorList variable
	prevSuccessorList := make([]string, SUCCESSOR_LIST_SIZE)
	copy(prevSuccessorList, n.GetSuccessorList())

	if hash.IsInRange(hash.Hash(x), hash.Hash(n.GetID()), hash.Hash(n.GetSuccessor())) {
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
	succSuccList, err := n.Gossiper.GetSuccessorList(n.GetID(), n.GetSuccessor())

	if err != nil {
		log.Printf("error in stabilize[GetSuccessorList]: %+v\n", err)
		n.fixSuccessorList()
		return
	}

	n.updateSuccessorList(succSuccList, prevSuccessorList)

	log.Println("Value of successorList: ", n.successorList)

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
	log.Printf("Fixing Successor List: %+v\n", n.GetSuccessorList())

	n.successorList = n.successorList[1:]
	n.successorList = append(n.successorList, "")

	isNewSuccessorAlive := n.healthCheck(n.GetSuccessor())
	if !isNewSuccessorAlive {
		n.fixSuccessorList()
	} else {

		if n.GetSuccessor() != "" {
			n.migrationFault(n.GetSuccessor())
		} else {
			// TODO: How to find the missing successor? Ring break
			return
		}

	}

}

func (n *Node) updateSuccessorList(succSuccList []string, prevSuccessorList []string) {
	copy(n.successorList[1:], succSuccList[:SUCCESSOR_LIST_SIZE-1])

	newElements, missingElements := compareList(prevSuccessorList, n.GetSuccessorList())

	if len(newElements) > 0 || len(missingElements) > 0 {
		keys, values := getAllLocalFiles()

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

func (n *Node) replicateToNodeList(nodeList []string, fileName, ip string) {
	for i := range nodeList {
		if nodeList[i] != n.GetID() {
			log.Printf("[REPLICATION] From Node %s to Node %s.\n", n.GetID(), nodeList[i])
			n.replicateToNode(nodeList[i], fileName, ip)
		}
	}
}

func (n *Node) deleteFromNodeList(nodeList []string, fileName string) {
	for i := range nodeList {
		if nodeList[i] != n.GetID() {
			log.Printf("[DELETE FROM REPLICA] From Node %s to Node %s.\n", n.GetID(), nodeList[i])
			_, err := n.Gossiper.DeleteFileFromNode(nodeList[i], fileName, "replica")
			if err != nil {
				print("Error in Deleting file: %+v\n", err)
			}
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
	// log.Printf("[NOTIFY HANDLER] Possible Predecessor of %s is %s.\n", n.GetID(), possiblePredecessor)
	//possiblePredecessor is Request's pred
	if (n.GetPredecessor() == "") ||
		(hash.IsInRange(
			hash.Hash(possiblePredecessor),
			hash.Hash(n.GetPredecessor()),
			hash.Hash(n.GetID()),
		)) {
		// log.Printf("[NOTIFY HANDLER - Set Predecessor] Predecessor of %s is %s.\n", n.GetID(), possiblePredecessor)
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
		log.Printf("%s's Predecessor is down.\n", n.GetID())
		n.SetPredecessor("")

	}
}
