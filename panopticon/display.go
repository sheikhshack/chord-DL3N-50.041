package main

import (
	"fmt"
	"sort"
	"time"
)

type nodeData struct {
	nodeID      string
	predecessor string
	successor   string
}

func (t Tower) display() {
	for {
		clearTerminal()
		fmt.Printf("%31s %15s %15s\n", "ID", "predecessor", "successor")
		var holder []nodeData
		t.data.Range(func(k interface{}, v interface{}) bool {
			node := v.(nodeData)
			holder = append(holder, node)
			//fmt.Printf("%15s: %15s %15s\n", id, node.predecessor, node.successor)
			return true
		})
		sort.SliceStable(holder, func(i, j int) bool {
			if holder[i].nodeID < holder[j].nodeID {
				return true
			}
			return false
		})
		for _, node := range holder {
			fmt.Printf("%15s: %15s %15s\n", node.nodeID, node.predecessor, node.successor)
		}
		//t.data.Range(func(k interface{}, v interface{}) bool {
		//	id := k.(string)
		//	node := v.(nodeData)
		//
		//	fmt.Printf("%15s: %15s %15s\n", id, node.predecessor, node.successor)
		//	return true
		//})
		time.Sleep(time.Millisecond * 500)
	}
}

func clearTerminal() {
	// clears the entire terminal
	fmt.Printf("\033[2J")
	// move cursor back to top left (aka print from the start again)
	fmt.Printf("\033[H")
}
