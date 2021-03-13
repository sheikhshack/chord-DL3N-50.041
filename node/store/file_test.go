package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	testString := "hello world"

	err := New("testwrite", []byte(testString))
	assert.Nilf(err, "error in creating file")

	filename, err := getFilename("testwrite")
	assert.Nilf(err, "error in obtaining filename")

	defer os.Remove(filename)

	content, err := ioutil.ReadFile(filename)
	assert.Nilf(err, "error in reading file")

	contentStr := string(content)
	assert.EqualValuesf(testString, contentStr, "different contents")
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	testString := "clown see clown follow"

	filename, err := getFilename("testread")
	assert.Nilf(err, "error in obtaining filename")

	err = ioutil.WriteFile(filename, []byte(testString), 0666)
	assert.Nilf(err, "error in creating file")

	defer os.Remove(filename)

	content, err := Get("testread")
	assert.Nilf(err, "error in reading file")

	contentStr := string(content)
	assert.EqualValuesf(testString, contentStr, "different contents")
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	testString := "clowns live in the circus"

	filename, err := getFilename("testdelete")
	assert.Nilf(err, "error in obtaining filename")
	err = ioutil.WriteFile(filename, []byte(testString), 0666)
	assert.Nilf(err, "error in creating file")

	err = Delete("testdelete")
	assert.Nilf(err, "error in removing file")
}
