package dl3n

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// DL3N represents the DL3N object asssociated with a particular file
// all fields except complete (and file, size in DL3NChunk) should have the same values for the same file
// struct should be marshallable to JSON (aka .dl3n file)
type DL3N struct {
	Mutex     *sync.Mutex `json:"-"`
	Name      string
	Hash      string
	Size      int64
	ChunkSize int64
	Chunks    []*DL3NChunk
}

// DL3N chunk represents a chunk of a DL3N object
type DL3NChunk struct {
	Id        int64
	Hash      string
	Size      int64
	Filepath  string `json:"-"` // nil if not available
	Available bool   `json:"-"`
}

// Create a new DL3N struct from filepath
// Handle chunking here
func NewDL3NFromFile(path string, chunkSize int64) (*DL3N, error) {
	// open the file once to get the infohash and fileSize
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	filename := filepath.Base(path)
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()

	infohash, err := getInfohash(f)
	if err != nil {
		return nil, err
	}
	f.Close()

	// open the file again to chunk
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	chunkCount, err := chunkFile(f, infohash, chunkSize)
	if err != nil {
		return nil, err
	}

	dl3n := DL3N{
		Mutex:     &sync.Mutex{},
		Name:      filename,
		Hash:      infohash,
		Size:      fileSize,
		ChunkSize: chunkSize,
		Chunks:    make([]*DL3NChunk, 0),
	}

	for i := int64(0); i < chunkCount; i++ {
		chunkPath := infohash + ".dl3nchunk." + strconv.FormatInt(i, 10)

		// open chunkfile to get hash
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return nil, err
		}
		chunkFileInfo, _ := chunkFile.Stat()
		chunkFileSize := chunkFileInfo.Size()

		chunkHash, err := getInfohash(chunkFile)
		if err != nil {
			return nil, err
		}
		chunkFile.Close()

		dl3n.Chunks = append(dl3n.Chunks, &DL3NChunk{
			Id:        i,
			Hash:      chunkHash,
			Size:      chunkFileSize,
			Filepath:  chunkPath,
			Available: true,
		})
	}

	return &dl3n, nil
}

// Create a new DL3N struct from filepath
// With one chunk
func NewDL3NFromFileOneChunk(path string) (*DL3N, error) {
	// open the file once to get the infohash and fileSize
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()
	f.Close()

	dl3n, err := NewDL3NFromFile(path, fileSize+1)

	return dl3n, err
}

// WriteMetaFile writes a .dl3n file to filepath for that particular DL3N
func (d *DL3N) WriteMetaFile(path string) error {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(d)

	if err != nil {
		return err
	}

	ioutil.WriteFile(path, buf.Bytes(), os.ModeAppend)

	return nil
}

// Create an empty DL3N struct from metadata filepath (.dl3n, actually just a JSON file)
func NewDL3NFromMeta(path string) (*DL3N, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := DL3N{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// gets SHA-1 infohash
func getInfohash(f *os.File) (string, error) {
	hash := sha1.New()

	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	infohash := fmt.Sprintf("%X", sum)

	return infohash, nil
}

// chunk
func chunkFile(f *os.File, infohash string, fileChunkSize int64) (int64, error) {
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()

	chunkCount := int64(math.Ceil(float64(fileSize) / float64(fileChunkSize)))

	for i := int64(0); i < chunkCount; i++ {

		partSize := int(math.Min(float64(fileChunkSize), float64(fileSize-int64(i*fileChunkSize))))
		partBuffer := make([]byte, partSize)

		f.Read(partBuffer)

		// write to disk
		fileName := infohash + ".dl3nchunk." + strconv.FormatInt(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			return 0, err
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
	}

	return chunkCount, nil
}

// unchunk
func joinChunks(chunks []*os.File, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, chunk := range chunks {
		b, err := ioutil.ReadAll(chunk)
		if err != nil {
			return err
		}
		outFile.Write(b)
	}

	return nil
}

// Complete returns if all DL3NChunks are available
func (d *DL3N) Complete() bool {
	for _, c := range d.Chunks {
		if !c.Available {
			return false
		}
	}

	return true
}
