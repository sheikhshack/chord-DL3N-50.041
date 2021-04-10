package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Sentry struct {
	ctx context.Context
	client *client.Client
	network string
}

func NewSentry(ctx context.Context, network string) *Sentry {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err!= nil {
		log.Fatal("Error connecting to Docker engine", err)
	}
	return &Sentry{ctx: ctx, client: client, network: network}
}

func (s *Sentry ) BuildChordImage()  {
	opt := types.ImageBuildOptions{
		Dockerfile:   "../Dockerfile",
		Tags:	[]string{"sheikhshack/ds-node"},
	}
	_, err := s.client.ImageBuild(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("Error building (latest) chord image, building from dockerhub instead, %v", err)
		out, err := s.client.ImagePull(s.ctx, "sheikhshack/ds-node", types.ImagePullOptions{})
		if err != nil {
			log.Fatalf("Failed all building routines: %v", err)
		}

		defer out.Close()
		io.Copy(os.Stdout, out)
	}
}
func (s *Sentry) SetupTestNetwork()  {
	// resets and remove current Network first
	_ = s.client.NetworkRemove(s.ctx, s.network)
	opt := types.NetworkCreate{Attachable: true}
	response, err := s.client.NetworkCreate(s.ctx, s.network, opt)
	if err != nil {
		log.Fatalf("Failed building network: %v", err)
	}
	log.Print(response)

}

func (s *Sentry ) FireOffChordNode(master bool, name string) {
	if master {

		configs := &container.Config{
			Hostname:        name,
			ExposedPorts:	 nat.PortSet{"9000/tcp": struct{}{}},
			Env:             []string{"PEER_HOSTNAME="},
			Image:           "sheikhshack/chord_node",
		}
		container, err := s.client.ContainerCreate(s.ctx, configs, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, name )
		if err != nil {
			log.Fatalf("Failed building container: %v", err)
		}
		fmt.Println(container)
		s.client.NetworkConnect(s.ctx, s.network, name, &network.EndpointSettings{})

	}
}



func main() {
	ctx := context.Background()
	sentry := NewSentry(ctx, "apache1")
	//sentry.BuildChordImage()
	sentry.SetupTestNetwork()
	sentry.FireOffChordNode(true, "cassandra-node")
}