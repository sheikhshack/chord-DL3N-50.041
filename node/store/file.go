package store

import (
	"io/ioutil"
	"os"
)

// New creates a new file locally with the nodeId directory, filename key and contains value
func New(nodeId, key string, value []byte) error {
	filename, err := getFilename(nodeId, key)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, value, 0666); err != nil {
		return err
	}

	return nil
}

// Get obtains the bytes stored in filename
func Get(nodeId, key string) ([]byte, error) {
	filename, err := getFilename(nodeId, key)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Get obtains the bytes stored in filename
func GetAll(nodeId string) ([]os.FileInfo, error) {
	dir, err := getOrCreateChordDir(nodeId)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Delete removes the file
func Delete(nodeId, key string) error {
	filename, err := getFilename(nodeId, key)
	if err != nil {
		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}
	return nil
}

// Migrate file from 1 folder to another
func Migrate(oldId, newId, key string) error {

	oldFileName, err := getFilename(oldId, key)
	if err != nil {
		return err
	}

	newFileName, err := getFilename(newId, key)
	if err != nil {
		return err
	}

	err = os.Rename(oldFileName, newFileName)
	if err != nil {
		return err
	}

	return nil
}
