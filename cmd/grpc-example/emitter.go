package main

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
	"log"
)

func main() {
	log.Printf("Starting Emitter!\n")
	c := chord.New("not relevant")

	res, err := c.Gossiper.Healthcheck(c.ID, "localhost")

	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}
	log.Printf("Got response!: %v", res)
}
