package store

import (
	"os"
	"path"

	"github.com/kardianos/osext"
)

// getExePath returns the path of the dir the executable is located at
func getExePath() (exPath string, err error) {
	return osext.ExecutableFolder()
}

// getOrCreateChordDir gets the path for chord store and creates it if it doesn't exist
func getOrCreateChordDir(exPath string) (dirPath string, err error) {
	dirPath = path.Join(exPath, "chord")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModeDir+0777)
		if err != nil {
			return "", err
		}
	}
	return dirPath, nil
}
