package client

import (
	"context"
	"fmt"
	"github.com/sheikhshack/distributed-chaos-50.041/basic"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// In this setup i think its not as efficient. Trigger only when needed, passive
type Node struct {
	id 				int
	hostname		string
	listenport		string
	sendport		string
	status			string

}

type RemoteNode struct {
	Hostname string
	Port     string
}

func New(id int, hostname, listenport, sendport, status string) *Node {
	return &Node{id, hostname, listenport, sendport, status}
}

// to be run as a goroutine....
func (n *Node) StartListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", n.listenport))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// proceeds to init the listener service provided by gprc + internal api
	s := basic.Listener{ID: n.id, Hostname: n.hostname}
	grpcServer := grpc.NewServer()
	basic.RegisterBasicServiceServer(grpcServer, &s)

	// throw error if serving fk up
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

func (n *Node) registerService() {
	//TODO: Have some stuff to run here first for service discovery

}

func (n *Node) SimPing(remote *Node) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":" + remote.listenport, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %s", err)

	}
	defer conn.Close()
	c := basic.NewBasicServiceClient(conn)
	message := basic.Message{Body: fmt.Sprintf("Ping from peer : %s", n.hostname)}
	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	log.Printf("Response from server: %s", response.Body)
}

func (n *Node) Ping(remote RemoteNode) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%s", remote.Hostname, remote.Port)
	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %s", err)

	}
	defer conn.Close()
	c := basic.NewBasicServiceClient(conn)
	message := basic.Message{Body: fmt.Sprintf("Ping from peer : %s", n.hostname)}
	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	log.Printf("Response from server: %s", response.Body)
}

func (n *Node) StartSim(remoteNode *Node) {
	go n.StartListener()

	for {
		time.Sleep(2 * time.Second)
		n.SimPing(remoteNode)
	}
}

func (n *Node) Start(remoteNode RemoteNode){
	go n.StartListener()

	for {
		time.Sleep(2 * time.Second)
		n.Ping(remoteNode)
	}
}