package main

import (
	"context"
	"github.com/sheikhshack/distributed-chaos-50.041/basic"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %s", err)

	}
	defer conn.Close()
	c := basic.NewBasicServiceClient(conn)
	message := basic.Message{Body: "Hello from the client"}

	response, err := c.SayHello(context.Background(), &message )
	if err != nil {
		log.Fatalf("Error sending message: %v", err)

	}
	log.Printf("Response from server: %s", response.Body)
}