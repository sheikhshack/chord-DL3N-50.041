package chord

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/log"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

const FINGER_TABLE_SIZE = 16

type Node struct {
	ID            string // maybe IP address
	fingers       []string
	predecessor   string
	next          int
	successorList []string
	replicaCount  int

	Gossiper *gossip.Gossiper
}

// New creates and returns a new Node
func New(id string) *Node {
	// 16 is finger table size
	replicaCount, err := strconv.Atoi(os.Getenv("SUCCESSOR_LIST_SIZE"))
	if err != nil {
		log.Error.Fatalf("Node INIT error, invalid SUCCESSOR_LIST_SIZE: %v", err)
	}
	n := &Node{ID: id, next: 0, fingers: make([]string, FINGER_TABLE_SIZE), successorList: make([]string, replicaCount), replicaCount: replicaCount}
	log.Info.Printf("Node:%v, HashedValue:%v\n", n.GetID(), hash.Hash(n.GetID()))
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
	files, err := store.GetAll("local")
	if err != nil {
		print(err)
	}

	for _, i := range files {
		log.Info.Printf("Filename:%v, HashedFile: %v", i.Name(), hash.Hash(i.Name()))
	}

	return n
}

func (n *Node) InitRing() {
	n.SetPredecessor("")
	n.SetSuccessor(n.GetID())
	go n.cron()
}

func (n *Node) Join(id string) {
	successor, err := n.Gossiper.Join(n.GetID(), id)
	if err != nil {
		// TODO: handle this error
		// we can pass the error back and have main.go to exit gracefully with helpful message
		log.Error.Fatalf("error in join: %+v\n", err)
	}
	n.SetPredecessor("")
	n.SetSuccessor(successor)
	go n.migrationJoin(n.GetSuccessor())
	//edge case of in the 1s window, the node's ideal pred hasn't recognised this node
	go n.cron()

}

func (n *Node) FindSuccessor(hashed int) string {
	if hash.IsInRange(hashed, hash.Hash(n.GetID()), hash.Hash(n.GetSuccessor())+1) {
		return n.GetSuccessor()
	} else {
		nPrime := n.closestPrecedingNode(hashed)

		successor, err := n.Gossiper.FindSuccessor(n.GetID(), nPrime, hashed)
		if err != nil {
			log.Warn.Printf("Finger %v is down: %+v\n", nPrime, err)

			for i := 0; i < FINGER_TABLE_SIZE; i++ {
				n.fixFingers()
			}
			time.Sleep(time.Millisecond * 200)
			// Assume that FingerTables will eventually be corrected if there is >1 node alive
			if n.successorList[0] != n.GetID() {
				return n.FindSuccessor(hashed)
			} else {
				// No more alive successor nodes except itself (Same as commented out edge case)
				return n.GetID()
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
		if hash.IsInRange(hash.Hash(n.fingers[i]), hash.Hash(n.GetID()), hashed) {
			return n.fingers[i]
		}
	}
	return n.GetID()
}

func (n *Node) WriteFiles(fileType, fileNames, ip string) error {
	key := fileNames
	val := ip

	var keys_list []string
	var val_list []string
	if strings.Contains(key, ",") {
		keys_list = strings.Split(key, ",")
		val_list = strings.Split(val, ",")
		keys_list = keys_list[:len(keys_list)-1]
		val_list = val_list[:len(val_list)-1]
	} else {
		keys_list = []string{key}
		val_list = []string{val}
	}

	for i := 0; i < len(keys_list); i++ {

		fileByte := []byte(val_list[i])
		output := store.New(fileType, keys_list[i], fileByte)
		if output != nil {
			log.Error.Printf("Error in writing to file:%v", output)
			return output
		}
	}

	return nil
}

func (n *Node) DeleteFiles(fileType, fileName string) error {
	key := fileName

	var keys_list []string

	if strings.Contains(key, ",") {
		keys_list = strings.Split(key, ",")
		keys_list = keys_list[:len(keys_list)-1]
	} else {
		keys_list = []string{key}
	}
	for i := 0; i < len(keys_list); i++ {
		output := store.Delete(fileType, keys_list[i])
		if output != nil {
			return output
		}
	}

	return nil
}

func (n *Node) WriteFileAndReplicate(fileType, fileName, ip string) error {
	err := n.WriteFiles(fileType, fileName, ip)
	if err != nil {
		return err
	}
	n.replicateToSuccessorList(fileName, ip)
	return nil
}

func (n *Node) DeleteFileAndReplicate(fileType, fileName string) error {
	err := n.DeleteFiles(fileType, fileName)
	if err != nil {
		return err
	}
	n.deleteFromSuccessorList(fileName)
	return nil
}
