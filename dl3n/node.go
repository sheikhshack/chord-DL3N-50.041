package dl3n

import (
	"context"
	"fmt"
	"net/http"
)

// Represents an application node
type DL3NNode struct {
	DL3N          *DL3N
	PeerDiscovery *PeerDiscovery
	server        *http.Server
}

type PeerDiscovery interface {
	FindPeer(string)
}

// NewDL3NNode
func NewDL3NNode(d *DL3N, p *PeerDiscovery) *DL3NNode {
	dn := &DL3NNode{
		DL3N:          d,
		PeerDiscovery: p,
	}

	return dn
}

// Basically starts a file server on addr
func (d *DL3NNode) StartSeed(addr string) {
	// TODO: allow other methods on server?
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, d.DL3N.Chunks[0].Filepath)
	}

	d.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(handler),
	}

	fmt.Printf("Starting Seed on %s ... \n", d.server.Addr)
	go d.server.ListenAndServe()
}

// Stops the file server
func (d *DL3NNode) StopSeed() {
	fmt.Printf("Stopping Seed on %s ... \n", d.server.Addr)
	d.server.Shutdown(context.Background())
}
