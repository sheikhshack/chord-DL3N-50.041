package store

import (
	"io/ioutil"
	"os"
)

// New creates a new file locally with the fileType directory (local / replica), filename key and contains value
func New(fileType, key string, value []byte) error {
	filename, err := getFilename(fileType, key)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, value, 0666); err != nil {
		return err
	}

	return nil
}

// Get obtains the bytes stored in filename with fileType string (local / replica) and key string
func Get(fileType, key string) ([]byte, error) {
	filename, err := getFilename(fileType, key)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Get obtains the bytes stored in filename with fileType string (local / replica)
func GetAll(fileType string) ([]os.FileInfo, error) {
	dir, err := getOrCreateChordDir(fileType)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Delete removes the file with fileType string (local / replica) and key string
func Delete(fileType, key string) error {
	filename, err := getFilename(fileType, key)
	if err != nil {
		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}
	return nil
}

// Migrate file from 1 folder to another
func Migrate(oldFileType, newFileType, key string) error {

	oldFileName, err := getFilename(oldFileType, key)
	if err != nil {
		return err
	}

	newFileName, err := getFilename(newFileType, key)
	if err != nil {
		return err
	}

	err = os.Rename(oldFileName, newFileName)
	if err != nil {
		return err
	}

	return nil
}
