package gossip

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/sheikhshack/distributed-chaos-50.041/log"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
)

// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

const (
	LISTEN_PORT     = 9000
	PANOPTICON_ADDR = "panopticon"
)

func (g *Gossiper) report() {
	if !g.DebugMode {
		return
	}

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", PANOPTICON_ADDR, LISTEN_PORT)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		//log.Fatalf("Cannot connect to: %s", err)
		return

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	_, err = client.Debug(context.Background(), &pb.DebugMessage{
		FromID:        g.Node.GetID(),
		Predecessor:   g.Node.GetPredecessor(),
		SuccessorList: g.Node.GetSuccessorList(),
		Fingers:       g.Node.GetFingers(),
	})
	if err != nil {
		//log.Printf("Error sending message: %v", err)
		return
	}
}

func (g *Gossiper) emit(nodeAddr string, request *pb.Request) (*pb.Response, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.Emit(context.Background(), request)
	if err != nil {
		//log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

// called by FindSuccessor
func (g *Gossiper) FindSuccessor(fromID, toID string, key int) (string, error) {
	req := &pb.Request{
		Command:     pb.Command_FIND_SUCCESSOR,
		RequesterID: fromID,
		TargetID:    toID,
		Body: &pb.Request_Body{
			HashSlot: int64(key),
		},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return "", err
	}
	return res.GetBody().GetID(), nil
}

// called by join
func (g *Gossiper) Join(fromID, toID string) (string, error) {
	//k = n.ID
	req := &pb.Request{
		Command:     pb.Command_JOIN,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return "", err
	}
	return res.GetBody().ID, nil
}

// called by checkPredecessor
// TODO: change all method signatures to include error in return
func (g *Gossiper) Healthcheck(fromID, toID string) (bool, error) {
	req := &pb.Request{
		Command:     pb.Command_HEALTHCHECK,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return false, err
	}
	return res.GetBody().IsHealthy, nil
}

//Get the predecessor of the node
func (g *Gossiper) GetPredecessor(fromID, toID string) (string, error) {
	req := &pb.Request{
		Command:     pb.Command_GET_PREDECESSOR,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return "", err
	}
	return res.GetBody().ID, nil
}

// Get the successor list of the node
func (g *Gossiper) GetSuccessorList(fromID, toID string) ([]string, error) {
	req := &pb.Request{
		Command:     pb.Command_GET_SUCCESSOR_LIST,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{},
	}

	res, err := g.emit(toID, req)
	if err != nil {
		return make([]string, 1), err
	}
	return res.GetBody().SuccessorList, nil
}

// called by notify
//n things it might be the predecessor of id
func (g *Gossiper) Notify(fromID, toID string) error {
	req := &pb.Request{
		Command:     pb.Command_NOTIFY,
		RequesterID: fromID,
		TargetID:    toID,
		Body:        &pb.Request_Body{},
	}

	_, err := g.emit(toID, req)
	if err != nil {
		return err
	}
	return nil
}

// external dialing service called w/o emit //
func (g *Gossiper) WriteFileToNode(nodeAddr, fileName, fileType, ip string) (*pb.ModResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)
	//log.Printf("Sending Request: %+v, %+v", request, request.Command)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := pb.NewInternalListenerClient(conn)
	response, err := client.WriteFile(context.Background(), &pb.ModRequest{
		Key:      fileName,
		Value:    ip,
		FileType: fileType,
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) WriteFileAndReplicateToNode(nodeAddr, fileName, fileType, ip string) (*pb.ModResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)
	//log.Printf("Sending Request: %+v, %+v", request, request.Command)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	log.Info.Printf("Writing file %v to Node %v to replicate\n", fileName, nodeAddr)
	client := pb.NewInternalListenerClient(conn)
	response, err := client.WriteFileAndReplicate(context.Background(), &pb.ModRequest{
		Key:      fileName,
		Value:    ip,
		FileType: fileType,
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) DeleteFileAndReplicateToNode(nodeAddr, fileName, fileType string) (*pb.ModResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)
	//log.Printf("Sending Request: %+v, %+v", request, request.Command)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	log.Info.Printf("Deleting file %v from Node %v to be replicated", fileName, nodeAddr)
	client := pb.NewInternalListenerClient(conn)
	response, err := client.DeleteFileAndReplicate(context.Background(), &pb.FetchChordRequest{
		Key:      fileName,
		FileType: fileType,
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) DeleteFileFromNode(nodeAddr, fileName, fileType string) (*pb.ModResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.DeleteFile(context.Background(), &pb.FetchChordRequest{
		Key:      fileName,
		FileType: fileType,
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) readFileFromNode(nodeAddr, fileName, fileType string) (*pb.ContainerInfo, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)
	//log.Printf("Sending Request: %+v, %+v", request, request.Command)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.ReadFile(context.Background(), &pb.FetchChordRequest{
		Key:      fileName,
		FileType: fileType,
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) MigrationJoinFromNode(nodeAddr string) (*pb.MigrationResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.MigrationJoin(context.Background(), &pb.MigrationRequest{
		RequesterID: g.Node.GetID(),
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}

func (g *Gossiper) MigrationFaultFromNode(nodeAddr string) (*pb.MigrationResponse, error) {
	g.report()

	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, LISTEN_PORT)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Error.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()

	client := pb.NewInternalListenerClient(conn)
	response, err := client.MigrationFault(context.Background(), &pb.MigrationRequest{
		RequesterID: g.Node.GetID(),
	})
	if err != nil {
		log.Warn.Printf("Error sending message: %v", err)
		return nil, err
	}
	return response, nil
}
