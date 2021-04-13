package gossip

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/sheikhshack/distributed-chaos-50.041/log"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

// methods available to the gossiper via the node package
type node interface {
	GetPredecessor() string
	GetSuccessorList() []string
	GetSuccessor() string
	SetSuccessor(id string)
	GetID() string
	GetFingers() []string

	FindSuccessor(hashed int) string
	NotifyHandler(possiblePredecessor string)
	MigrationJoinHandler(requestId string)
	MigrationFaultHandler(requestId string)
	WriteFiles(fileType, fileName, ip string) error
	WriteFileAndReplicate(fileType, fileName, ip string) error
	DeleteFiles(fileType, fileName string) error
}

type Gossiper struct {
	Node      node
	DebugMode bool

	pb.UnimplementedInternalListenerServer
}

func (g *Gossiper) NewServerAndListen(listenPort int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		log.Error.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInternalListenerServer(s, g)
	if err := s.Serve(lis); err != nil {
		log.Error.Fatalf("failed to serve: %v", err)
	} else {
		log.Info.Printf("Listening on port %v\n", listenPort)
	}
	return s
}

// This method is server-side
// TODO: consider that we need to spin up goroutine whenever we serve a new request
func (g *Gossiper) Emit(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	g.report()
	var res *pb.Response

	switch in.Command {
	case pb.Command_FIND_SUCCESSOR:
		key := int(in.GetBody().HashSlot)
		id := g.findSuccessorHandler(key)
		res = &pb.Response{
			Command:     pb.Command_FIND_SUCCESSOR,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: id,
			},
		}

	case pb.Command_JOIN:
		fromID := in.GetRequesterID()
		id := g.joinHandler(fromID)
		res = &pb.Response{
			Command:     pb.Command_JOIN,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: id,
			},
		}

	case pb.Command_HEALTHCHECK:
		isHealthy := g.healthcheckHandler()
		res = &pb.Response{
			Command:     pb.Command_HEALTHCHECK,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				IsHealthy: isHealthy,
			},
		}

	case pb.Command_GET_PREDECESSOR:
		predecessorID := g.getPredecessorHandler()
		res = &pb.Response{
			Command:     pb.Command_GET_PREDECESSOR,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				ID: predecessorID,
			},
		}

	case pb.Command_GET_SUCCESSOR_LIST:
		successorList := g.getSuccessorListHandler()
		res = &pb.Response{
			Command:     pb.Command_GET_SUCCESSOR_LIST,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body: &pb.Response_Body{
				SuccessorList: successorList,
			},
		}

	case pb.Command_NOTIFY:
		possiblePredecessor := in.GetRequesterID()
		g.notifyHandler(possiblePredecessor)
		res = &pb.Response{
			Command:     pb.Command_NOTIFY,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body:        &pb.Response_Body{},
		}

	case pb.Command_REPLICATE_TO_NODE:

		key := in.GetBody().Key
		value := in.GetBody().Value
		fileType := in.GetBody().FileType

		err := g.replicateToNodeHandler(fileType, key, value)
		if err != nil {
			return nil, err
		}

		res = &pb.Response{
			Command:     pb.Command_REPLICATE_TO_NODE,
			RequesterID: in.GetRequesterID(),
			TargetID:    in.GetTargetID(),
			Body:        &pb.Response_Body{},
		}

	default:
		return nil, errors.New("command not recognised")
	}

	return res, nil
}

// handler to findSuccessor
func (g *Gossiper) findSuccessorHandler(key int) (id string) {
	return g.Node.FindSuccessor(key)
}

// handler to join
func (g *Gossiper) joinHandler(fromID string) string {
	log.Info.Printf("Node %v request to join ring\n", fromID)
	//fromID= previous node's id
	if g.Node.GetID() == g.Node.GetSuccessor() {
		g.Node.SetSuccessor(fromID)
		return g.Node.GetID()
	}
	return g.Node.FindSuccessor(hash.Hash(fromID))
}

// handler to healthcheck (checkPredecessor)
func (g *Gossiper) healthcheckHandler() bool {
	// can return false if Node deems itself unhealthy
	return true
}

func (g *Gossiper) getPredecessorHandler() string {
	return g.Node.GetPredecessor()
}

func (g *Gossiper) getSuccessorListHandler() []string {
	return g.Node.GetSuccessorList()
}

// notifyHandler handles notify requests and returns if id is in between n.predecessor and n.
// notifyHandler might also update n.predecessor and trigger data transfer if appropriate.
func (g *Gossiper) notifyHandler(possiblePredecessor string) {
	g.Node.NotifyHandler(possiblePredecessor)
}

func (g *Gossiper) replicateToNodeHandler(fileType, key, value string) error {
	return g.Node.WriteFiles(fileType, key, value)
}

////////// Internal CHORD file management
// FetchChordIp is basically unused for now, leaving as legacy
func (g *Gossiper) FetchChordIp(ctx context.Context, fetchRequest *pb.FetchChordRequest) (*pb.ModResponse, error) {
	key := fetchRequest.GetKey()
	ip := g.Node.FindSuccessor(hash.Hash(key))
	return &pb.ModResponse{IP: ip}, nil
}

func (g *Gossiper) WriteFile(ctx context.Context, writeRequest *pb.ModRequest) (*pb.ModResponse, error) {

	output := g.Node.WriteFiles(writeRequest.FileType, writeRequest.Key, writeRequest.Value)
	return &pb.ModResponse{IP: g.Node.GetID()}, output
}

func (g *Gossiper) DeleteFile(ctx context.Context, fetchRequest *pb.FetchChordRequest) (*pb.ModResponse, error) {
	output := g.Node.DeleteFiles(fetchRequest.FileType, fetchRequest.GetKey())
	return &pb.ModResponse{IP: g.Node.GetID()}, output
}

// ReadFile allows returning of container IP from file within chord
func (g *Gossiper) ReadFile(ctx context.Context, fetchRequest *pb.FetchChordRequest) (*pb.ContainerInfo, error) {
	key := fetchRequest.GetKey()
	fileType := fetchRequest.FileType

	fileByte, status := store.Get(fileType, key)
	// TODO: Might need to change this
	containerIP := string(fileByte[:])
	//log.Printf("--- FS: Triggering File Read in Node for key [%v] \n", key)

	return &pb.ContainerInfo{ContainerIP: containerIP}, status
}

func (g *Gossiper) MigrationJoin(ctx context.Context, migrationRequest *pb.MigrationRequest) (*pb.MigrationResponse, error) {
	requestId := migrationRequest.RequesterID
	//log.Printf("--- FS: Triggering Migration (Join) to Chord Node %v from %v \n", requestId, g.Node.GetID())

	g.Node.MigrationJoinHandler(requestId)

	return &pb.MigrationResponse{Success: true}, nil
}

func (g *Gossiper) MigrationFault(ctx context.Context, migrationRequest *pb.MigrationRequest) (*pb.MigrationResponse, error) {
	requestId := migrationRequest.RequesterID
	//log.Printf("--- FS: Triggering Migration (Fault) to Chord Node %v from %v \n", requestId, g.Node.GetID())

	g.Node.MigrationFaultHandler(requestId)

	return &pb.MigrationResponse{Success: true}, nil
}

func (g *Gossiper) WriteFileAndReplicate(ctx context.Context, writeRequest *pb.ModRequest) (*pb.ModResponse, error) {

	output := g.Node.WriteFileAndReplicate(writeRequest.FileType, writeRequest.Key, writeRequest.Value)
	return &pb.ModResponse{IP: g.Node.GetID()}, output
}

/////////////////////////////////////////
////////// EXPOSED TO D3LOWEN ///////////
/////////////////////////////////////////

// StoreKeyHash will allow  D3L to call a single function to write upload directory to the right chord node
func (g *Gossiper) StoreKeyHash(ctx context.Context, dlUploadRequest *pb.DLUploadRequest) (*pb.DLResponse, error) {
	fileName := dlUploadRequest.Filename
	containerIP := dlUploadRequest.ContainerIP
	log.Info.Printf("[EXT] Storing file [%v] from [%v]\n", fileName, containerIP)

	correctChordIP := g.Node.FindSuccessor(hash.Hash(fileName))
	_, err := g.WriteFileAndReplicateToNode(correctChordIP, fileName, "local", containerIP)
	if err != nil {
		// TODO: return error
		log.Error.Fatalf("error in running Store Key Hash (ext) : %+v\n", err)
	}
	status := pb.DlowenStatus_SAME_NODE
	if correctChordIP != g.Node.GetID() {
		status = pb.DlowenStatus_REDIRECTED_NODE
	}

	return &pb.DLResponse{
		ChordIP: correctChordIP,
		Status:  status,
	}, nil

}

// GetFileLocation takes in a file and resolves it into an IP
func (g *Gossiper) GetFileLocation(ctx context.Context, dlDownloadRequest *pb.DLDownloadRequest) (*pb.DLDownloadResponse, error) {
	fileName := dlDownloadRequest.Filename
	correctChordIP := g.Node.FindSuccessor(hash.Hash(fileName))
	containerInfo, err := g.readFileFromNode(correctChordIP, fileName, "local")
	if err != nil {
		// TODO: return error
		log.Error.Fatalf("error in running GetFile Location (ext) : %+v\n", err)
	}
	return &pb.DLDownloadResponse{Container: containerInfo, ChordIP: correctChordIP}, nil

}
