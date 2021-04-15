package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"log"
	"os"
)

func main() {
	var currHostname string
	manualDNS := os.Getenv("MY_PEER_DNS")
	if manualDNS == "DEFAULT" {
		id, err := os.Hostname()
		if err != nil {
			log.Fatal("Docker engine -- hostname issues", err)
		}
		currHostname = id
	} else {
		currHostname = manualDNS
	}

	knownPeerID := os.Getenv("PEER_HOSTNAME")
	log.Printf("Starting NODE_ID: %v\n", currHostname)

	node := chord.New(currHostname)

	if knownPeerID == "" {
		node.InitRing()
		log.Printf("%v: Ring setup\n", node.ID)
	} else {
		node.Join(knownPeerID)
		log.Printf("%v: Joined ring\n", node.ID)
	}

	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
}
