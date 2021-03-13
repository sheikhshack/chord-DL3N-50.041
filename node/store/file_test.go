package store

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	testString := "hello world"

	err := New("testwrite", []byte(testString))
	assert.Nilf(err, "error in creating file")

	exPath, err := getExePath()
	assert.Nilf(err, "error in getting executable path")

	filename := path.Join(exPath, "chord", "testwrite")

	defer os.Remove(filename)

	content, err := ioutil.ReadFile(filename)
	assert.Nilf(err, "error in reading file")

	contentStr := string(content)
	assert.EqualValuesf(testString, contentStr, "different contents")
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	testString := "clown see clown follow"

	exPath, err := getExePath()
	assert.Nilf(err, "error in getting executable path")
	dirPath, err := getOrCreateChordDir(exPath)
	assert.Nilf(err, "error in creating directory")

	filename := path.Join(dirPath, "testread")

	err = ioutil.WriteFile(filename, []byte(testString), 0666)
	assert.Nilf(err, "error in creating file")

	defer os.Remove(filename)

	content, err := Get("testread")
	assert.Nilf(err, "error in reading file")

	contentStr := string(content)
	assert.EqualValuesf(testString, contentStr, "different contents")
}
