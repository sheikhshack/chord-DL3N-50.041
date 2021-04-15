package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	sentryFS "github.com/sheikhshack/distributed-chaos-50.041/integration-tests/sentry"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func getHostVol(containerName string) string {
	relativePath := fmt.Sprintf("./volumes/%s/", containerName)
	dir, err := filepath.Abs(filepath.Dir(relativePath))
	if err != nil {
		panic(err)
	}
	return dir
}

type Sentry struct {
	ctx          context.Context
	client       *client.Client
	network      string
	replicaCount int
	master       string
	slaves       []string
}

func NewSentry(ctx context.Context, network string, replicaCount int, master string, slaves []string) *Sentry {
	log.Println("--Sentry: Setting up docker environment")
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error connecting to Docker engine", err)
	}
	// remove all previous settings and running containers
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		client.ContainerKill(ctx, container.ID, "SIGKILL")
		client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})

	}

	volumes, err := client.VolumeList(ctx, filters.Args{})
	if err != nil {
		panic(err)
	}
	for _, vol := range volumes.Volumes {
		if err := client.VolumeRemove(ctx, vol.Name, true); err != nil {
			panic(err)
		}

	}

	client.ContainersPrune(ctx, filters.Args{})
	client.NetworksPrune(ctx, filters.Args{})
	client.VolumesPrune(ctx, filters.Args{}) // please make sure you have no volume from other stuff before running

	cmd := exec.Command("sudo", "rm", "-rf", "./volumes")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return &Sentry{ctx: ctx, client: client, network: network, replicaCount: replicaCount, master: master, slaves: slaves}
}

// Legacy: BuildChordImage should be commented out unless you know what youre doing
func (s *Sentry) BuildChordImage() {
	opt := types.ImageBuildOptions{
		Dockerfile: "../Dockerfile",
		Tags:       []string{"sheikhshack/chord_node"},
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

func (s *Sentry) FindContainerID(name string) string{
	containers, err := s.client.ContainerList(s.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if container.Names[0] == "/"+name {
			return container.ID
		}

	}
	return ""
}

func (s *Sentry) FindNetworkID(name string) string{
	containers, err := s.client.NetworkList(s.ctx, types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}
	for _, net := range containers {
		if net.Name == name{
			return net.ID
		}
	}
	return ""
}
// Legacy: BuildMikeImage should be commented out, unless you know what youre doing
func (s *Sentry) BuildMikeImage() {
	opt := types.ImageBuildOptions{
		Dockerfile: "../Dockerfile.mike",
		Tags:       []string{"sheikhshack/mike_node"},
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

// SetupTestNetwork sets up the common docker network
func (s *Sentry) SetupTestNetwork() {
	log.Println("--Sentry: Setting up chord network")

	// resets and remove current Network first
	if err := s.client.NetworkRemove(s.ctx, s.network); err != nil {
		//log.Printf("Failed to remove network, %v\n", err)
	}
	opt := types.NetworkCreate{Attachable: true}
	_, err := s.client.NetworkCreate(s.ctx, s.network, opt)
	if err != nil {
		log.Fatalf("Failed building network: %v", err)
	}


}

// FireOffMikeNode starts the mike client (for testing purposes only)
func (s *Sentry) FireOffMikeNode(contactNode, name, cmd1, cmd2 string) {

	s.client.ContainerRemove(s.ctx, name, types.ContainerRemoveOptions{Force: true})
	attachedNode := fmt.Sprintf("APP_NODE=%v", contactNode)

	env := []string{attachedNode}

	configs := &container.Config{
		Hostname: name,
		Env:      env,
		Image:    "sheikhshack/mike_node_tester",
		Cmd:      []string{"writefile", cmd1, cmd2},
	}
	container, err := s.client.ContainerCreate(s.ctx, configs, nil, nil, nil, name)
	if err != nil {
		log.Fatalf("Failed building container: %v", err)
	}
	s.client.NetworkConnect(s.ctx, s.network, name, &network.EndpointSettings{})
	if err := s.client.ContainerStart(s.ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to run container: %v", err)
	}
}

// FireOffChordNode basically starts the container
func (s *Sentry) FireOffChordNode(ringLeader bool, name string) {

	s.client.ContainerStop(s.ctx, name, nil)
	s.client.ContainerRemove(s.ctx, name, types.ContainerRemoveOptions{Force: true})
	replicaConfig := fmt.Sprintf("SUCCESSOR_LIST_SIZE=%v", s.replicaCount)
	dnsName := fmt.Sprintf("MY_PEER_DNS=%s", name)

	var env []string
	if ringLeader {
		name = "alpha"
		dnsName = fmt.Sprintf("MY_PEER_DNS=%s", name)
		env = []string{"PEER_HOSTNAME=", replicaConfig, dnsName}

	} else {
		env = []string{"PEER_HOSTNAME=alpha", replicaConfig, dnsName}
	}

	// Provisions volume mount point first
	if err := os.MkdirAll("./volumes/"+name, os.ModePerm); err != nil {
		panic(err)
	}

	configs := &container.Config{
		Hostname:     name,
		ExposedPorts: nat.PortSet{"9000/tcp": struct{}{}},
		Env:          env,
		Image:        "sheikhshack/chord_node",
	}
	bindingHostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:         mount.TypeBind,
				Source:       getHostVol(name),
				Target:       "/built-app/chord",
				ReadOnly:     false,
				TmpfsOptions: &mount.TmpfsOptions{Mode: os.ModePerm},
			},
		},
	}
	container, err := s.client.ContainerCreate(s.ctx, configs, bindingHostConfig, nil, nil, name)
	if err != nil {
		log.Fatalf("Failed building container: %v", err)
	}
	s.client.NetworkConnect(s.ctx, s.network, name, &network.EndpointSettings{})
	if err := s.client.ContainerStart(s.ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to run container: %v", err)
	}
}

func (s *Sentry) StopContainer(name string) {
	containerID := s.FindContainerID(name)
	s.client.ContainerKill(s.ctx, containerID,"SIGKILL" )

}

func (s *Sentry) ForceStopContainer(name string) {
	containerID := s.FindContainerID(name)
	s.client.ContainerKill(s.ctx, containerID, "SIGKILL")
	s.client.ContainerRemove(s.ctx, containerID, types.ContainerRemoveOptions{Force: true})

}

func (s *Sentry) StartContainerAgain(name string)  {
	containerID := s.FindContainerID(name)
	s.client.ContainerRestart(s.ctx, containerID,nil )

}

// Writes a file via a directed chord node
func (s *Sentry) WriteFileToChord(viaNode, fileName, content string) {
	log.Printf("--Sentry: Writing to ringFS for fileName %s\n", fileName)

	command1 := fmt.Sprintf("-f=%s", fileName)
	command2 := fmt.Sprintf("-c=%s", content)
	s.FireOffMikeNode(viaNode, "mike_test", command1, command2)
}

func (s *Sentry) InterruptConnection (name, networkName string) {
	netID := s.FindNetworkID(networkName)
	cont := s.FindContainerID(name)
	s.client.NetworkDisconnect(s.ctx, netID, cont, true)

}

func (s *Sentry) ReestablishConnection (name, networkName string) {
	netID := s.FindNetworkID(networkName)
	cont := s.FindContainerID(name)
	s.client.NetworkConnect(s.ctx, netID, cont, nil)

}

//func (s *Sentry) RestartContainer(name string) {
//	contID := s.FindContainerID(name)
//
//	dur:= 50 * time.Second
//	fmt.Println("-- MECH: Restarting node", contID)
//	s.client.ContainerRestart(s.ctx, contID, &dur)
//}

// Procedure to bring up the test case ring
func (s *Sentry) BringUpChordRing() {
	log.Println("--Sentry: Setting up RING")

	s.SetupTestNetwork()
	// fireoff the master first, then the slaves
	s.FireOffChordNode(true, s.master)
	for _, v := range s.slaves {
		go s.FireOffChordNode(false, v)
	}

}

func (s *Sentry) BringDownRing() {
	log.Println("--Sentry: Bringing down docker environment")

	containers, err := s.client.ContainerList(s.ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		s.client.ContainerKill(s.ctx, container.ID, "SIGKILL")
		s.client.ContainerRemove(s.ctx, container.ID, types.ContainerRemoveOptions{Force: true})

	}

	volumes, err := s.client.VolumeList(s.ctx, filters.Args{})
	if err != nil {
		panic(err)
	}
	for _, vol := range volumes.Volumes {
		if err := s.client.VolumeRemove(s.ctx, vol.Name, true); err != nil {
			panic(err)
		}
	}

	s.client.ContainersPrune(s.ctx, filters.Args{})
	s.client.NetworksPrune(s.ctx, filters.Args{})
	s.client.VolumesPrune(s.ctx, filters.Args{}) // please make sure you have no volume from other stuff before running

	cmd := exec.Command("sudo", "rm", "-rf", "./volumes")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Sentry) ReportChordFS() {

}

func test2() {
	ctx := context.Background()
	//INIT Test case 1 - Replication Setup of 3
	testReplica := 3
	master := "alpha"
	slaves := []string{"nodeBravo", "nodeDelta"}
	sentry := NewSentry(ctx, "test2-network", testReplica, master, slaves)
	sentry.BringUpChordRing()
	time.Sleep(time.Second * 10)
	fmt.Println("-- TEST2.2: Sending the files over to a node (random) w/ system replica set to ", testReplica)

	sentry.WriteFileToChord("alpha", "alpha", "Alpha file content")
	sentry.WriteFileToChord("alpha", "nodeBravo", "Bravo file content")
	sentry.WriteFileToChord("alpha", "nodeCharlie", "Charlie file content")
	sentry.WriteFileToChord("alpha", "nodeDelta", "Delta file content")
	time.Sleep(time.Second * 10)
	sentryFS.ReadFileInVolume()
	fmt.Println("-- TEST2.3: Bringing up Node charlie (previously not existent) ")
	sentry.FireOffChordNode(false, "nodeCharlie")
	time.Sleep(15 * time.Second)
	fmt.Println("-- TEST2.3: Current chord file system as such:")
	sentryFS.ReadFileInVolume()
	fmt.Println("-- TEST2.4: Bringing down chord node charlie")
	sentry.ForceStopContainer("nodeCharlie")
	sentryFS.DeleteFilesystemLink("nodeCharlie")
	time.Sleep(5 * time.Second)
	fmt.Println("-- TEST2.4: Current chord file system as such:")
	sentryFS.ReadFileInVolume()
	sentry.BringDownRing()

}

//INIT Test case 1 - Replication Setup of 3
func test1() {
	ctx := context.Background()
	testReplica := 3
	master := "apache"
	slaves := []string{"slave-node1", "slave-node2", "slave-node3", "slave-node4", "slave-node5"}
	sentry := NewSentry(ctx, "apache1", testReplica, master, slaves)
	sentry.BringUpChordRing()
	time.Sleep(time.Second * 15)
	sentry.WriteFileToChord("slave-node1", "slave-node1", "I love 50.041: Distributed Systems and Computing!")
	fmt.Println("-- TEST1: Sending the file over to a node (random) w/ system replica set to ", testReplica)
	time.Sleep(time.Second * 10)
	sentryFS.ReadFileInVolume()
	fmt.Println("-- TEST1: End of procedure ")

	current_replica := "slave-node1"
	sentry.ForceStopContainer(current_replica)
	fmt.Println("-- TEST2: Removing current primary", current_replica)
	sentryFS.DeleteFilesystemLink(current_replica)
	time.Sleep(time.Second * 5)
	sentryFS.ReadFileInVolume()
	sentry.BringDownRing()

	// INIT Test case 2 -
}

func test3() {
	ctx := context.Background()
	testReplica := 3
	master := "apache"
	slaves := []string{"nodeBravo", "nodeCharlie", "nodeDelta", "nodeGamma", "node5"}
	sentry := NewSentry(ctx, "apache1", testReplica, master, slaves)
	sentry.BringUpChordRing()
	time.Sleep(time.Second * 15)
	sentry.WriteFileToChord("nodeBravo", "nodeBravo", ":(")
	fmt.Println("-- TEST3: Sending the file over to a node (random) w/ system replica set to ", testReplica)
	time.Sleep(time.Second * 10)
	sentryFS.ReadFileInVolume()

	current_replica := "nodeBravo"
	successive_replica := "nodeDelta"
	sentry.ForceStopContainer(current_replica)
	sentry.ForceStopContainer(successive_replica)
	fmt.Println("-- TEST3: Removing two consecutive nodes ", current_replica, successive_replica)
	sentryFS.DeleteFilesystemLink(current_replica)
	sentryFS.DeleteFilesystemLink(successive_replica)
	time.Sleep(time.Second * 20)
	sentryFS.ReadFileInVolume()
}

func test4 () {
	ctx := context.Background()
	network := "apache1"
	testReplica := 3
	master := "apache"
	slaves := []string{"slave-node1", "slave-node2", "slave-node3", "slave-node4", "slave-node5"}
	sentry := NewSentry(ctx, network, testReplica, master, slaves)

	sentry.BringUpChordRing()
	time.Sleep(time.Second *  15)
	fmt.Println("-- TEST4: Bringing slave-node1 out of network")
	sentry.WriteFileToChord("slave-node1", "slave-node1", "I love 50.041: Distributed Systems and Computing!")
	time.Sleep(time.Second * 10)
	fmt.Println("-- TEST4: Current file system before fault")
	sentryFS.ReadFileInVolume()

	sentry.StopContainer("slave-node1")
	time.Sleep(time.Second * 5)
	fmt.Println("-- TEST4: Current file system for reads during loss of 1")
	sentryFS.ReadFileInVolume()

	sentry.StartContainerAgain("slave-node1")
	time.Sleep(time.Second * 5)
	fmt.Println("-- TEST4: Current file system after 1 joins back the ring")
	sentryFS.ReadFileInVolume()

}



func main() {
	//test1()
	test2()
	//test3()
	//test4()
}
