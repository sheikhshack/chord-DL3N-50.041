package store

import (
	"github.com/kardianos/osext"
	"os"
	"path"
)

// getOrCreateChordDir gets the path for chord store and creates it if it doesn't exist
func getOrCreateChordDir() (dirPath string, err error) {
	exPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	dirPath = path.Join(exPath, "chord")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModeDir+0777)
		if err != nil {
			return "", err
		}
	}
	return dirPath, nil
}

func getFilename(name string) (filename string, err error) {
	if dirPath, err := getOrCreateChordDir(); err != nil {
		return "", err
	} else {
		filename = path.Join(dirPath, name)
		return filename, nil
	}
}
