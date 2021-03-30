package dl3n

import (
	"context"
	"fmt"
	"os"

	pb "github.com/sheikhshack/distributed-chaos-50.041/dl3n/proto"
	"google.golang.org/grpc"
)

type DL3NNodeClient struct {
	PeerDiscoveryService PeerDiscoveryService
	DL3N                 DL3N
}

func (d *DL3NNodeClient) Get(stop chan bool) error {
	// check path for complete chunks
	for _, c := range d.DL3N.Chunks {
		filename := fmt.Sprintf("%s.dl3nchunk.%d", c.Hash, c.Id)
		f, err := os.Open(filename)
		if err != nil {
			continue
		}

		fileHash, err := getInfohash(f)
		if err != nil {
			continue
		}
		f.Close()

		if fileHash != c.Hash {
			os.Remove(filename)
		}

		c.Available = true
		c.Filepath = fileHash
	}

	for !d.Complete() {
		for _, c := range d.DL3N.Chunks {
			if c.Available {
				continue
			}

			for _, p := range d.PeerDiscoveryService.FindPeers(d.DL3N.Hash) {

			}

		}
	}

	return nil
}

func getChunkData(chunkId int64, infohash, addr string) ([]byte, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewDL3NClient(conn)
	ctx := context.Background()
	r, err := c.GetChunk(ctx, &pb.ChunkMeta{
		DL3NHash: &pb.DL3NHash{Hash: infohash},
		ChunkId:  chunkId,
	})

	if err != nil {
		return nil, err
	}

	b := r.Data

	return b, nil
}
