package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
	"time"
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
	client.ContainersPrune(ctx, filters.Args{})
	client.NetworksPrune(ctx, filters.Args{})
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
	if err := s.client.NetworkRemove(s.ctx, s.network); err != nil {
		log.Printf("Failed to remove network, %v\n", err)
	}
	opt := types.NetworkCreate{Attachable: true}
	response, err := s.client.NetworkCreate(s.ctx, s.network, opt)
	if err != nil {
		log.Fatalf("Failed building network: %v", err)
	}
	log.Print(response)

}

func (s *Sentry ) FireOffMikeNode(contactNode, name, cmd1, cmd2 string) {
	s.client.ContainerRemove(s.ctx, name , types.ContainerRemoveOptions{Force: true})
	attachedNode := fmt.Sprintf("APP_NODE=%v", contactNode)
	env := []string{attachedNode}

	configs := &container.Config{
		Hostname:        name,
		Env:             env ,
		Image:           "sheikhshack/mike_node_tester",
		Cmd:   			 []string{"writefile", cmd1, cmd2},
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
	s.client.ContainerStop(s.ctx, name, nil)
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

func (s * Sentry) WriteFileToChord (viaNode, fileName, content string) {
	command1 := fmt.Sprintf("-f %s", fileName)
	command2 := fmt.Sprintf("-c %s", content)
	s.FireOffMikeNode(viaNode, "mike_test", command1, command2)


}



func main() {
	ctx := context.Background()
	sentry := NewSentry(ctx, "apache1")
	////sentry.BuildChordImage()
	sentry.SetupTestNetwork()
	sentry.FireOffChordNode(true, "master-node")
	sentry.FireOffChordNode(false, "slave-node1")
	sentry.FireOffChordNode(false, "slave-node2")
	sentry.FireOffChordNode(false, "slave-node3")


	time.Sleep(time.Second *  20)
	sentry.WriteFileToChord("slave-node1", "wombat.txt", "I hate this shit")

}