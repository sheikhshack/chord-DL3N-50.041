package gossip

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
)

// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

const (
	LISTEN_PORT = 9000
	EMIT_PORT   = 9001
)

func dial(nodeAddr string) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%s", nodeAddr, LISTEN_PORT)
	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
	log.Printf("Response from server: %s", response.Body)
}

// called by FindSuccessor
func FindSuccessor(fromID, toID string, key int) string {
	panic("not implemented")
}

// called by join
func Join(fromID, toID string) string {
	//k = n.ID
	panic("not implemented")
}

// called by checkPredecessor
func Healthcheck(fromID, toID string) bool {
	panic("not implemented")
}

//Get the predecessor of the node
func GetPredecessor(fromID, toID string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the predecessor of id
func Notify(fromID, toID string) {
	//pred = n.ID
	panic("not implemented")
}

// Not used?
// Called by Lookup
// TODO: move this method to exposed API package
func Get(key string) ([]byte, error) {
	panic("not implemented")
}
