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
		d, err := dl3n.NewDL3NFromFile(path, 64)
		fmt.Print(d)
		fmt.Print(err)
	}

}
