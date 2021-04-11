package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
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
		Tags:	[]string{"sheikhshack/chord_node"},
	}
	_, err := s.client.ImageBuild(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("Error building (latest) chord image, building from dockerhub instead, %v", err)
		out, err := s.client.ImagePull(s.ctx, "sheikhshack/chord_node", types.ImagePullOptions{})
		if err != nil {
			log.Fatalf("Failed all building routines: %v", err)
		}

		defer out.Close()
		io.Copy(os.Stdout, out)
	}
}

func (s *Sentry ) BuildMikeImage()  {
	opt := types.ImageBuildOptions{
		Dockerfile:   "../Dockerfile.mike",
		Tags:	[]string{"sheikhshack/mike_node"},
	}
	_, err := s.client.ImageBuild(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("Error building (latest) mike image, building from dockerhub instead, %v", err)
		out, err := s.client.ImagePull(s.ctx, "sheikhshack/mike_node", types.ImagePullOptions{})
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

func (s *Sentry ) FireOffMikeNode(contact string, name string) {
	s.client.ContainerRemove(s.ctx, name , types.ContainerRemoveOptions{Force: true})
	attachedNode := fmt.Sprintf("APP_NODE=%v", name)
	env := []string{attachedNode}

	configs := &container.Config{
		Hostname:        name,
		ExposedPorts:	 nat.PortSet{"9000/tcp": struct{}{}},
		Env:             env ,
		Image:           "sheikhshack/chord_node",
	}
	container, err := s.client.ContainerCreate(s.ctx, configs, nil, nil,  nil, name )
	if err != nil {
		log.Fatalf("Failed building container: %v", err)
	}
	fmt.Println(container)
	s.client.NetworkConnect(s.ctx, s.network, name, &network.EndpointSettings{})
	if err:= s.client.ContainerStart(s.ctx, container.ID, types.ContainerStartOptions{}); err!= nil{
		log.Fatalf("Failed to run container: %v", err)
	}
}

func (s *Sentry ) FireOffChordNode(ringLeader bool, name string) {
	s.client.ContainerRemove(s.ctx, name , types.ContainerRemoveOptions{Force: true})
	var env []string
	if ringLeader {
		env = []string{"PEER_HOSTNAME="}
		name = "alpha"
	} else {
		env= []string{"PEER_HOSTNAME=alpha"}
	}


	configs := &container.Config{
		Hostname:        name,
		ExposedPorts:	 nat.PortSet{"9000/tcp": struct{}{}},
		Env:             env ,
		Image:           "sheikhshack/chord_node",
	}
	container, err := s.client.ContainerCreate(s.ctx, configs, nil, nil,  nil, name )
	if err != nil {
		log.Fatalf("Failed building container: %v", err)
	}
	fmt.Println(container)
	s.client.NetworkConnect(s.ctx, s.network, name, &network.EndpointSettings{})
	if err:= s.client.ContainerStart(s.ctx, container.ID, types.ContainerStartOptions{}); err!= nil{
		log.Fatalf("Failed to run container: %v", err)
	}
}

func (s * Sentry) CheckFile (fileName string) {

}



func main() {
	ctx := context.Background()
	sentry := NewSentry(ctx, "apache1")
	//sentry.BuildChordImage()
	sentry.SetupTestNetwork()
	sentry.FireOffChordNode(true, "master-node")
	sentry.FireOffChordNode(false, "slave-node")

}