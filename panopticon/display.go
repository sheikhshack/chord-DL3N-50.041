package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
)

type nodeData struct {
	nodeID      string
	predecessor string
	successor   string
}

func (t Tower) display() {
	for {
		clearTerminal()
		fmt.Printf("%39s %23s %23s %10s\n", "ID", "predecessor", "successor", "stabilized")
		var holder []nodeData
		t.data.Range(func(k interface{}, v interface{}) bool {
			node := v.(nodeData)
			holder = append(holder, node)
			return true
		})

		sort.SliceStable(holder, func(i, j int) bool {
			if holder[i].nodeID < holder[j].nodeID {
				return true
			}
			return false
		})

		for _, node := range holder {
			hID := hash.Hash(node.nodeID)
			hpredecessor := hash.Hash(node.predecessor)
			hsuccessor := hash.Hash(node.successor)
			stabilized := node.predecessor != "" && node.successor != "" && hash.IsInRange(hID, hpredecessor, hsuccessor)

			fmt.Printf("%15s (%v): %15s (%v) %15s (%v) %9v\n",
				node.nodeID, hID,
				node.predecessor, hpredecessor,
				node.successor, hsuccessor,
				stabilized,
			)
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func clearTerminal() {
	// clears the entire terminal
	fmt.Printf("\033[2J")
	// move cursor back to top left (aka print from the start again)
	fmt.Printf("\033[H")
}
