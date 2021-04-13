package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"log"
	"os"
)

func main() {

	id, err := os.Hostname()
	if err != nil {
		log.Fatal("Docker engine -- hostname issues", err)
	}
	knownPeerID := os.Getenv("PEER_HOSTNAME")
	log.Printf("Starting NODE_ID: %v\n", id)

	node := chord.New(id)

	if knownPeerID == "" {
		node.InitRing()
		log.Printf("%v: Ring setup\n", node.ID)
	} else {
		node.Join(knownPeerID)
		log.Printf("%v: Joined ring\n", node.ID)
	}

	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
}
