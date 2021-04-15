package sentryFS

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func ReadFileInVolume() {
	fmt.Println("--- Current FS State -----------")
	printDirectory("volumes", "./volumes", 0)

	//err := filepath.Walk("./volumes",
	//	func(path string, info os.FileInfo, err error) error {
	//		if err != nil {
	//			return err
	//		}
	//		if info.IsDir() == false{
	//			filer, _ := os.Open(path)
	//			b, _:= ioutil.ReadAll(filer)
	//			fmt.Printf("%s : %s\n", path, b)
	//		}
	//		return nil
	//	})
	//if err != nil {
	//	log.Println(err)
	//}
}

func ReadFiles() {
	fileInfos, err := ioutil.ReadDir("./volumes")
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}
	for _, file := range fileInfos {
		fmt.Println(file.Name())
	}
}

func DeleteFilesystemLink(name string) {
	fmt.Println("-- FS: Removing FS for ", name)
	cmd := exec.Command("sudo", "rm", "-rf", "./volumes"+"/"+name)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
