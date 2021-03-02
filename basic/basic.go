package basic

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
)

type Listener struct {
	 ID			int
	 Hostname 	string
}

func (s *Listener) SayHello (ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &Message{Body: fmt.Sprintf("ACK from node %s : %d", s.Hostname, s.ID)}, nil
}