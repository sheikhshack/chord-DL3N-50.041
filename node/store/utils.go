package store

import (
	"errors"
	"os"
	"path"

	"github.com/kardianos/osext"
)

// getOrCreateChordDir gets the path for chord store and creates it if it doesn't exist
func getOrCreateChordDir(fileType string) (dirPath string, err error) {
	exPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	switch fileType {

	case "local":
		dirPath = path.Join(exPath, "chord", fileType)

	case "replica":
		dirPath = path.Join(exPath, "chord", fileType)

	default:
		err := errors.New("invalid file type")
		return "", err

	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModeDir+0777)
		if err != nil {
			return "", err
		}
	}
	return dirPath, nil
}

func getFilename(fileType, name string) (filename string, err error) {
	if dirPath, err := getOrCreateChordDir(fileType); err != nil {
		return "", err
	} else {
		filename = path.Join(dirPath, name)
		return filename, nil
	}
}
