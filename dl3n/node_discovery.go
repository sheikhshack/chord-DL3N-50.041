package dl3n

import (
	"context"
	"fmt"

	gossip "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"google.golang.org/grpc"
)

type ChordNodeDiscovery struct {
	nodeAddr string
}

func NewChordNodeDiscovery(nodeAddr string) *ChordNodeDiscovery {
	return &ChordNodeDiscovery{
		nodeAddr: nodeAddr,
	}
}

// SetSeederAddr is basically mike's writeExternalFile
func (nd *ChordNodeDiscovery) SetSeederAddr(infohash, addr string) error {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nd.nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	_, err = client.StoreKeyHash(context.Background(), &gossip.DLUploadRequest{
		Filename:    infohash,
		ContainerIP: addr,
	})
	if err != nil {
		return err
	}

	return nil
}

// FindSeederAddr is basically mike's resolveFile
func (nd *ChordNodeDiscovery) FindSeederAddr(infohash string) (string, error) {
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nd.nodeAddr, 9000)
	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.GetFileLocation(context.Background(), &gossip.DLDownloadRequest{
		Filename: infohash,
	})
	if err != nil {
		return "", err
	}

	ip := res.GetContainer().GetContainerIP()

	return ip, nil
}
