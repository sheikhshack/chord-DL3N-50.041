package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

func main() {
	helpMessage := "Usage - dl3n [seed|get] [filepath] [addr]"

	if len(os.Args) != 4 {
		fmt.Println(helpMessage)
	}

	cmd := os.Args[1]
	path := os.Args[2]
	addr := os.Args[3]

	if cmd != "seed" && cmd != "get" {
		fmt.Println(helpMessage)
	}

	if cmd == "create" {
		d, _ := dl3n.NewDL3NFromFile(path, 64)
		d.WriteMetaFile(path + ".dl3n")
	}

	if cmd == "get" {
		d, _ := dl3n.NewDL3NFromMeta(path)
		fmt.Printf("%+v\n", d)
	}

	if cmd == "seed" {
		d, err := dl3n.NewDL3NFromFileOneChunk(path)
		if err != nil {
			fmt.Print(err)
		}

		d.WriteMetaFile(path + ".dl3n")

		ds := dl3n.NewDL3NNode(d, nil)
		ds.StartSeed(addr)

		// listen for signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		// wait for sigint or sigterm
		<-sigs
		ds.StopSeed()
	}

}
