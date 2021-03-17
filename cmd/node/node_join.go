package main

import (
	"fmt"
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/utils"
	"log"
	"os"
	"time"
)

func main() {
	_ = os.Mkdir("tmp", os.ModeDir+0777)
	logFileName := fmt.Sprintf("./tmp/log_%.23s.txt", time.Now().UTC())

	utils.SetLogFile(logFileName)

	id := os.Getenv("NODE_ID")
	knownPeerID := os.Getenv("PEER_HOSTNAME")
	log.Printf("starting NODE_ID: %v\n", id)

	node := chord.New(id)

	if knownPeerID == "" {
		node.InitRing()
		log.Printf("%v: init-ed ring\n", node.ID)
	} else {
		node.Join(knownPeerID)
		log.Printf("%v: joined peer %v\n", node.ID, knownPeerID)
	}

	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
}
