package chord

import (
	"github.com/sheikhshack/distributed-chaos-50.041/log"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
	"strings"
)

func (n *Node) replicateToNode(toID, fileName, ip string) {
	_, err := n.Gossiper.WriteFileToNode(toID, fileName, "replica", ip)
	if err != nil {
		log.Warn.Printf("Error in instructing node %v to replicate file: %+v\n", toID, ip)
	}
}

func (n *Node) replicateToSuccessorList(fileName, ip string) {

	// Assumes that successorList nodes repeat after it contains own node
	for i := range n.successorList {
		if n.successorList[i] != n.GetID() {
			n.replicateToNode(n.successorList[i], fileName, ip)
		} else {
			break
		}
	}
}

func (n *Node) replicateToNodeList(nodeList []string, fileName, ip string) {
	for i := range nodeList {
		if nodeList[i] != n.GetID() {
			n.replicateToNode(nodeList[i], fileName, ip)
		}
	}
}

func (n *Node) deleteFromNode(toID, fileName string) {
	_, err := n.Gossiper.DeleteFileFromNode(toID, fileName, "replica")
	if err != nil {
		log.Warn.Printf("Error in instructing node %v to deleting file: %+v\n", toID, err)
	}
}

func (n *Node) deleteFromNodeList(nodeList []string, fileName string) {
	for i := range nodeList {
		if nodeList[i] != n.GetID() {
			n.deleteFromNode(nodeList[i], fileName)
		}
	}
}

func (n *Node) deleteFromSuccessorList(fileName string) {
	// Assumes that successorList nodes repeat after it contains own node
	for i := range n.successorList {
		if n.successorList[i] != n.GetID() {
			n.deleteFromNode(n.successorList[i], fileName)
		} else {
			break
		}
	}
}

//ask successor to migrate files that belong to current node
func (n *Node) migrationJoin(successor string) {
	if _, err := n.Gossiper.MigrationJoinFromNode(successor); err != nil {
		log.Error.Fatalf("Error in migrating files from successor: %+v\n", err)
	}
}

//predecessor has asked to migrate files from current nodes
func (n *Node) MigrationJoinHandler(requestID string) {

	// Get all the replica files in the store
	files, err := store.GetAll("local")
	if err != nil {
		log.Error.Printf("error in reading file system: %+v\n", err)
		return
	}
	keys := ""
	values := ""
	//Loop through them and write over the ones that do not lie in between pred and current node, and then delete if the write is successful
	for _, i := range files {
		if !hash.IsInRange(hash.Hash(i.Name()), hash.Hash(requestID)+1, hash.Hash(n.GetID())+1) {
			keys += i.Name() + ","

			val, _ := store.Get("local", i.Name())
			values += string(val) + ","
		}
	}
	if keys != "" {
		_, err = n.Gossiper.WriteFileToNode(requestID, keys, "local", values)
		if err != nil {
			log.Error.Printf("Error in writing file: %+v\n", err)
			return
		} else {
			keysList := strings.Split(keys, ",")

			for _, i := range keysList {
				store.LocalMigrate("local", "replica", i)
			}
		}

		//Init deleting from last successor
		_, err = n.Gossiper.DeleteFileFromNode(n.successorList[n.replicaCount-1], keys, "replica")
		if err != nil {
			log.Error.Printf("Error in Deleting file: %+v\n", err)
			return
		}
	}

}

//ask successor to migrate files that belong to current node
func (n *Node) migrationFault(livenode string) {
	if _, err := n.Gossiper.MigrationFaultFromNode(livenode); err != nil {
		log.Error.Fatalf("Error in Migrating files from Successor: %+v\n", err)
	}
}

//predecessor has asked to migrate files from current nodes
func (n *Node) MigrationFaultHandler(requestID string) {

	// Get all the replica files in the store
	files, err := store.GetAll("replica")
	if err != nil {
		log.Error.Printf("error in reading file system: %+v\n", err)
		return
	}
	keys := ""
	values := ""
	//Loop through them and write over the ones that do not lie in between pred and current node, and then delete if the write is successful
	for _, i := range files {
		//log.Info.Printf("Has Filename:%v, HashedFile: %v, Checking if in between %v and %v", i.Name(), hash.Hash(i.Name()), hash.Hash(requestID), hash.Hash(n.GetID()))
		if hash.IsInRange(hash.Hash(i.Name()), hash.Hash(requestID), hash.Hash(n.GetID())) {
			keys += i.Name() + ","
			val, _ := store.Get("replica", i.Name())
			values += string(val) + ","
		}
	}
	if keys != "" {
		//Init writing to predecessor
		//log.Info.Printf("Migrating Filename:%v", keys)
		n.replicateToSuccessorList(keys, values)
		keysList := strings.Split(keys, ",")
		for _, i := range keysList {
			store.LocalMigrate("replica", "local", i)

		}
	}
}
