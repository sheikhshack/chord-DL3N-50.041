package main

import (
	"fmt"
	"github.com/sheikhshack/distributed-chaos-50.041/client"
)

func main() {

	node0 := client.New(123, "Alpha", "9001", "9002", "Alive")
	node1 := client.New(789, "Bravo", "9004", "9005", "Alive")
	go node0.StartSim(node1)
	go node1.StartSim(node0)
	fmt.Scanf("Enter")

}