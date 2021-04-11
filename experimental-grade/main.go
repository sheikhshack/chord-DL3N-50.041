package main

import (
	"log"
	"os"
	"time"
)

func main()  {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal("Host failed to get")
	}
	for {
		log.Print(host)
		time.Sleep(time.Second * 1)
	}
}