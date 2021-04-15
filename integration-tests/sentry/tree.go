package sentryFS

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func printListing(entry string, depth int) {
	indent := strings.Repeat("|   ", depth)
	fmt.Printf("%s|-- %s\n", indent, entry)
}

func printDirectory(dirName string, dirPath string, depth int) {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", dirPath, err.Error())
		return
	}

	printListing(dirName, depth)
	for _, entry := range entries {
		if (entry.Mode() & os.ModeSymlink) == os.ModeSymlink {
			full_path, err := os.Readlink(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				fmt.Printf("error reading link: %s\n", err.Error())
			} else {
				printListing(entry.Name()+" -> "+full_path, depth+1)
			}
		} else if entry.IsDir() {
			printDirectory(entry.Name(), filepath.Join(dirPath, entry.Name()), depth+1)
		} else {
			printListing(entry.Name(), depth+1)
		}
	}
}
