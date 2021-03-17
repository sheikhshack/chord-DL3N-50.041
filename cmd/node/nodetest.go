package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"log"
)

func master() {
	node.InitRing()
	log.Printf("%v: init-ed ring\n", node.ID)
	log.Printf("%v's predecessor: %s \n", node.ID, node.GetPredecessor)
	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
	log.Printf("%v's predecessor: %s \n", node.ID, node.GetPredecessor)
}

func worker() {
	node.Join(knownPeerID)
	log.Printf("%v: joined peer %v\n", node.ID, knownPeerID)
	log.Printf("%v's predecessor: %s \n", node.ID, node.GetPredecessor)
	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
	log.Printf("%v's predecessor: %s \n", node.ID, node.GetPredecessor)
}
