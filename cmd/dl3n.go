package main

import (
	"fmt"
	"os"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

func main() {
	helpMessage := "Usage - dl3n [seed|get] [filepath] [addr]"

	if len(os.Args) != 4 {
		fmt.Println(helpMessage)
	}

	cmd := os.Args[1]
	path := os.Args[2]
	// addr := os.Args[3]

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
		// d, _ := dl3n.NewDL3NFromFile(path, 64)
		// d.WriteMetaFile(path + ".dl3n")

		d, err := dl3n.NewDL3NFromFileOneChunk(path)
		if err != nil {
			fmt.Print(err)
		}

		d.WriteMetaFile(path + ".dl3n")

		// s := dl3n.DL3NNodeServer{
		// 	DL3N: *d,
		// 	Addr: addr,
		// }

		// interrupt := make(chan os.Signal, 1)
		// signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		// s.Seed(interrupt)
	}

}
