package sentryFS

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadFileInVolume(){
	fmt.Println("--- Current FS State -----------")
	err := filepath.Walk("./volumes",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false{
				filer, _ := os.Open(path)
				b, _:= ioutil.ReadAll(filer)
				fmt.Printf("%s : %s\n", path, b)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func ReadFiles(){
	fileInfos, err := ioutil.ReadDir("./volumes")
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}
	for _, file := range fileInfos {
		fmt.Println(file.Name())
	}
}