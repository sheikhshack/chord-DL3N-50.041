package dl3n

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Represents an application node
type DL3NNode struct {
	DL3N          *DL3N
	NodeDiscovery NodeDiscovery
	server        *http.Server
}

type NodeDiscovery interface {
	FindSeederAddr(string) string
}

// NewDL3NNode
func NewDL3NNode(d *DL3N, s NodeDiscovery) *DL3NNode {
	dn := &DL3NNode{
		DL3N:          d,
		NodeDiscovery: s,
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

func (d *DL3NNode) Get() error {
	nd := d.NodeDiscovery
	addr := nd.FindSeederAddr(d.DL3N.Hash)

	// Get the data
	resp, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(d.DL3N.Name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
