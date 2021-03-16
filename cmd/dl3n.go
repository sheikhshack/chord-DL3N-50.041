package main

import (
	"fmt"
	"os"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

func main() {
	helpMessage := "Usage - dl3n [create|seed|get] [filepath]"

	if len(os.Args) != 3 {
		fmt.Println(helpMessage)
	}

	cmd := os.Args[1]
	path := os.Args[2]

	if cmd != "create" && cmd != "seed" && cmd != "get" {
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

}
