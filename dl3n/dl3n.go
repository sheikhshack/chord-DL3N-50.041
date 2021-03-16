package dl3n

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

// DL3N represents the DL3N object asssociated with a particular file
// all fields except complete (and file, size in DL3NChunk) should have the same values for the same file
// struct should be marshallable to JSON (aka .dl3n file)
type DL3N struct {
	Hash      string
	Size      int64
	ChunkSize int64
	Chunks    []*DL3NChunk
	Complete  bool `json:"-"` // true if all chunks are available
}

// DL3N chunk represents a chunk of a DL3N object
type DL3NChunk struct {
	Id        int64
	Hash      string
	Size      int64
	File      *os.File `json:"-"` // nil if not available
	Available bool     `json:"-"`
}

// Create a new DL3N struct from filepath
// Handle chunking here
func NewDL3NFromFile(path string, chunkSize int64) (*DL3N, error) {
	// open the file once to get the infohash and fileSize
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
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
		Hash:      infohash,
		Size:      fileSize,
		ChunkSize: chunkSize,
		Chunks:    make([]*DL3NChunk, 0),
		Complete:  true,
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

		// open chunkfile again to attach to dl3nchunk
		chunkFile, err = os.Open(chunkPath)
		if err != nil {
			return nil, err
		}

		dl3n.Chunks = append(dl3n.Chunks, &DL3NChunk{
			Id:        i,
			Hash:      chunkHash,
			Size:      chunkFileSize,
			File:      chunkFile,
			Available: true,
		})
	}

	return &dl3n, nil
}

// WriteMetaFile writes a .dl3n file to filepath for that particular DL3N
func (d *DL3N) WriteMetaFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, b, os.ModeAppend)

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
