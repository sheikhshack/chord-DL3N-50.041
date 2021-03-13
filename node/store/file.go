package store

import (
	"io/ioutil"
	"path"
)

// New creates a new file locally with the filename key and contains value
func New(key string, value []byte) error {
	exPath, err := getExePath()
	if err != nil {
		return err
	}

	dirPath, err := getOrCreateChordDir(exPath)
	if err != nil {
		return err
	}

	filename := path.Join(dirPath, key)

	err = ioutil.WriteFile(filename, value, 0666)
	if err != nil {
		return err
	}

	return nil
}

// Get obtains the bytes stored in filename
func Get(key string) ([]byte, error) {
	exPath, err := getExePath()
	if err != nil {
		return nil, err
	}

	dirPath, err := getOrCreateChordDir(exPath)
	if err != nil {
		return nil, err
	}

	filename := path.Join(dirPath, key)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}
