package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"log"
)

func main() {
	log.Printf("Starting Listener!!\n")
	server := chord.New("localhost")
	server.InitRing()

	server.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
}
