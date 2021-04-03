package store

import (
	"io/fs"
	"io/ioutil"
	"os"
)

// New creates a new file locally with the filename key and contains value
func New(key string, value []byte) error {
	filename, err := getFilename(key)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, value, 0666); err != nil {
		return err
	}

	return nil
}

// Get obtains the bytes stored in filename
func Get(key string) ([]byte, error) {
	filename, err := getFilename(key)
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
func GetAll() ([]fs.FileInfo, error) {
	dir, err := getOrCreateChordDir()
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
func Delete(key string) error {
	filename, err := getFilename(key)
	if err != nil {
		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}
	return nil
}
