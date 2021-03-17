package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"github.com/sheikhshack/distributed-chaos-50.041/node/gossip"
	"github.com/sheikhshack/distributed-chaos-50.041/node/utils"
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
		time.Sleep(time.Millisecond * 5000)
		node.Join(knownPeerID)

	}

	node.Gossiper.NewServerAndListen(gossip.LISTEN_PORT)
}
