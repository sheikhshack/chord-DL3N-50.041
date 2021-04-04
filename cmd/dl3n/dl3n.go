package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

const helpMessage = `
DL3N SEED - usage:
dl3n seed [filepath] [addr]
filepath: file to be seeded. should be in same directory as executable.
addr: address of chord node to update for peer discovery. should include port.

DL3N GET - usage:
dl3n get [filepath] [addr]
filepath: .dl3n meta file. should be in same directory as executable.
addr: address of chord node to query for peer discovery. should include port.
`

func main() {
	if len(os.Args) != 4 {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	path := os.Args[2]
	addr := os.Args[3]

	if cmd != "seed" && cmd != "get" {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if cmd == "get" {
		d, _ := dl3n.NewDL3NFromMeta(path)
		fmt.Printf("%+v\n", d)

		nd := dl3n.NewChordNodeDiscovery(addr)

		ds := dl3n.NewDL3NNode(d, nd)
		err := ds.Get()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if cmd == "seed" {
		d, err := dl3n.NewDL3NFromFileOneChunk(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		d.WriteMetaFile(path + ".dl3n")

		nd := dl3n.NewChordNodeDiscovery(addr)
		ds := dl3n.NewDL3NNode(d, nd)

		containerIP := os.Getenv("CONTAINER_IP")
		seederAddr := containerIP + ":11111"
		err = ds.NodeDiscovery.SetSeederAddr(d.Hash, seederAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ds.StartSeed(seederAddr)

		// listen for signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		// wait for sigint or sigterm
		<-sigs
		ds.StopSeed()
	}

}
