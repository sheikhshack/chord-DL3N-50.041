package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

type ServerState struct {
	Mutex    *sync.Mutex `json:"-"`
	State    string      // "WAIT" or "SEEDING" or "GETTING" or "GET_DONE"
	SeedMeta *dl3n.DL3N
	GetMeta  *dl3n.DL3N
}

type GuiServer struct {
	NodeDiscovery dl3n.NodeDiscovery
	seederAddr    string
	seederNode    *dl3n.DL3NNode
	getterNode    *dl3n.DL3NNode
}

func NewGuiServer(nd dl3n.NodeDiscovery, s string) *GuiServer {
	g := &GuiServer{
		NodeDiscovery: nd,
		seederAddr:    s,
	}

	return g
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (g *GuiServer) StartServer() {
	serverState := &ServerState{
		Mutex: &sync.Mutex{},
		State: "WAIT",
	}

	getStateHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		serverState.Mutex.Lock()

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)

		serverState.Mutex.Unlock()
	}

	uploadHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		r.ParseMultipartForm(10 << 20)
		inFile, inHandler, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer inFile.Close()

		fPath := "./" + inHandler.Filename
		f, err := os.Create(fPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, inFile)

		d, _ := dl3n.NewDL3NFromFileOneChunk(fPath)

		serverState.Mutex.Lock()

		serverState.SeedMeta = d

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)

		serverState.Mutex.Unlock()
	}

	startSeedHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		serverState.Mutex.Lock()
		defer serverState.Mutex.Unlock()
		g.seederNode = dl3n.NewDL3NNode(serverState.SeedMeta, g.NodeDiscovery)
		serverState.State = "SEEDING"

		g.NodeDiscovery.SetSeederAddr(serverState.SeedMeta.Hash, g.seederAddr)
		g.seederNode.StartSeed(g.seederAddr)

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)
	}

	stopSeedHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		serverState.Mutex.Lock()
		defer serverState.Mutex.Unlock()

		g.seederNode.StopSeed()
		g.seederNode = nil
		serverState.State = "WAIT"
		serverState.SeedMeta = nil

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)
	}

	uploadMetaHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		r.ParseMultipartForm(10 << 20)
		inFile, _, err := r.FormFile("uploadMeta")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer inFile.Close()

		fPath := "./meta.dl3n"
		f, err := os.Create(fPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, inFile)

		d, _ := dl3n.NewDL3NFromMeta(fPath)

		serverState.Mutex.Lock()

		serverState.GetMeta = d

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)

		serverState.Mutex.Unlock()
	}

	startGetHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		serverState.Mutex.Lock()
		g.getterNode = dl3n.NewDL3NNode(serverState.GetMeta, g.NodeDiscovery)
		serverState.State = "GETTING"

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(serverState)
		serverState.Mutex.Unlock()

		go func() {
			g.getterNode.Get()

			serverState.Mutex.Lock()
			serverState.State = "GET_DONE"
			serverState.Mutex.Unlock()
		}()
	}

	getFileHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		serverState.Mutex.Lock()
		filePath := serverState.GetMeta.Name
		serverState.Mutex.Unlock()

		http.ServeFile(w, r, filePath)

		serverState.Mutex.Lock()
		serverState.State = "WAIT"
		serverState.Mutex.Unlock()
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		io.WriteString(w, index)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/getState", getStateHandler)
	http.HandleFunc("/startSeed", startSeedHandler)
	http.HandleFunc("/stopSeed", stopSeedHandler)
	http.HandleFunc("/uploadMeta", uploadMetaHandler)
	http.HandleFunc("/startGet", startGetHandler)
	http.HandleFunc("/getFile", getFileHandler)
	http.HandleFunc("/", indexHandler)

	addr := "0.0.0.0:3000"
	fmt.Printf("Starting gui server on %s ... \n", addr)
	http.ListenAndServe(addr, nil)
}
