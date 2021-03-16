package dl3n

import (
	"crypto/sha1"
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
	Size      int
	ChunkSize int
	Chunks    []*DL3NChunk
	Complete  bool // true if all chunks are available
}

// DL3N chunk represents a chunk of a DL3N object
type DL3NChunk struct {
	Id        int
	Hash      string
	Size      int
	File      *os.File `json:"-"` // nil if not available
	Available bool
}

// Create a new DL3N struct from filepath
// Handle chunking here
func NewDL3NFromFile(path string, chunkSize int) (*DL3N, error) {
	return nil, nil
}

// Create an empty DL3N struct from metadata filepath (.dl3n, actually just a JSON file)
func NewDL3NFromMeta(path string) (*DL3N, error) {
	return nil, nil
}

// WriteMetaFile writes a .dl3n file to filepath for that particular DL3N
func (*DL3N) WriteMetaFile(path string) error {
	return nil
}

// gets SHA-1 infohash
func getInfohash(f *os.File) (string, error) {
	hash := sha1.New()

	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	infohash := string(sum)

	return infohash, nil
}

// chunk
func chunkFile(f *os.File, infohash string, chunkSize int) error {
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()

	fileChunkSize := int64(chunkSize * (1 << 10)) // fileChunkSize is in bytes, chunkSize is in kbytes
	chunkCount := int64(math.Ceil(float64(fileSize) / float64(fileChunkSize)))

	for i := int64(0); i < chunkCount; i++ {

		partSize := int(math.Min(float64(fileChunkSize), float64(fileSize-int64(i*fileChunkSize))))
		partBuffer := make([]byte, partSize)

		f.Read(partBuffer)

		// write to disk
		fileName := infohash + ".dl3nchunk." + strconv.FormatInt(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			return err
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
	}

	return nil
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
