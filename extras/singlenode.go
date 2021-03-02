package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/client"
	"log"
	"os"
	"strconv"
)

// default configs. Our protocol will use ports 9001 and 9002 by defauly
const listenPort = "9001"
const sendPort = "9002"

// Setup the nodes accordingly

func handleErrors(err error){
	if err != nil {
		log.Fatal("Failed with error: ", err)
	}
}

func main () {

	hostname, err := os.Hostname()
	handleErrors(err)
	// get the dynamic ENV vars set in docker
	id := os.Getenv("NODE_ID")
	idString, err := strconv.Atoi(id)
	handleErrors(err)

	// get the other vars
	samplePeer := client.RemoteNode{
		Hostname: os.Getenv("PEER_HOSTNAME"),
		Port: os.Getenv("PEER_PORT"),
	}

	SERVICE := client.New(idString, hostname, listenPort, sendPort, "ALIVE")
	SERVICE.Start(samplePeer)

	// scanf for troubleshootings



}
