package store

import (
	"os"
	"path"

	"github.com/kardianos/osext"
)

// getOrCreateChordDir gets the path for chord store and creates it if it doesn't exist
func getOrCreateChordDir(nodeId string) (dirPath string, err error) {
	exPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	dirPath = path.Join(exPath, "chord", nodeId)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModeDir+0777)
		if err != nil {
			return "", err
		}
	}
	return dirPath, nil
}

func getFilename(nodeId, name string) (filename string, err error) {
	if dirPath, err := getOrCreateChordDir(nodeId); err != nil {
		return "", err
	} else {
		filename = path.Join(dirPath, name)
		return filename, nil
	}
}
