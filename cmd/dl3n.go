package main

import (
	"fmt"
	"os"
)

func main() {
	helpMessage := "Usage - dl3n [create|seed|get] [filepath]"

	if len(os.Args) != 3 {
		fmt.Println(helpMessage)
	}

	cmd := os.Args[1]
	// path := os.Args[2]

	if cmd != "create" && cmd != "seed" && cmd != "get" {
		fmt.Println(helpMessage)
	}
}
