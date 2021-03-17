package dl3n

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	pb "github.com/sheikhshack/distributed-chaos-50.041/dl3n/proto"
	"google.golang.org/grpc"
)

type DL3NNodeServer struct {
	pb.UnimplementedDL3NServer
	DL3N DL3N
	Addr string
}

func (s *DL3NNodeServer) GetAvailableChunks(ctx context.Context, h *pb.DL3NHash) (*pb.Chunks, error) {
	s.DL3N.Mutex.Lock()
	defer s.DL3N.Mutex.Unlock()

	if h.Hash != s.DL3N.Hash {
		return nil, errors.New("server is not seeding this dl3n")
	}

	chunks := pb.Chunks{
		Chunks: make([]*pb.ChunkMeta, 0),
	}

	for _, c := range s.DL3N.Chunks {
		if c.Available {
			chunks.Chunks = append(chunks.Chunks, &pb.ChunkMeta{
				ChunkId:  c.Id,
				DL3NHash: &pb.DL3NHash{Hash: c.Hash},
			})
		}
	}

	return &chunks, nil
}

func (s *DL3NNodeServer) GetChunk(ctx context.Context, m *pb.ChunkMeta) (*pb.ChunkData, error) {
	s.DL3N.Mutex.Lock()
	defer s.DL3N.Mutex.Unlock()

	for _, c := range s.DL3N.Chunks {
		if m.DL3NHash.Hash == s.DL3N.Hash && m.ChunkId == c.Id && c.Available {
			d, err := ioutil.ReadFile(c.Filepath)
			if err != nil {
				return nil, err
			}

			chunkData := pb.ChunkData{
				Data: d,
			}

			return &chunkData, nil
		}
	}

	return nil, errors.New("server is not seeding this chunk")
}

func (s *DL3NNodeServer) Seed(stop chan os.Signal) error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	svr := grpc.NewServer()
	pb.RegisterDL3NServer(svr, &DL3NNodeServer{})

	fmt.Printf("Seeding DL3N at %s ...", s.Addr)

	go svr.Serve(lis)
	<-stop
	svr.GracefulStop()

	return nil
}
