package gui

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

// temporarily using a string for dev purposes
// will switch to using a static file later

const indexHtml = `
<html>
    <body>
        <form
            enctype="multipart/form-data"
            action="/upload"
            method="post"
        >
            <input type="file" name="myFile" />
            <input type="submit" value="upload" />
        </form>
    </body>
</html>
`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexHtml)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("/upload")

	r.ParseMultipartForm(10 << 20)

	inFile, inHandler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	inFile.Close()

	fPath := "./" + inHandler.Filename
	f, err := os.Create(fPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Close()

	io.Copy(f, inFile)

	d, _ := dl3n.NewDL3NFromFileOneChunk(fPath)

	d.WriteMetaFile("./metafile.dl3n")

	metaF, _ := os.Open("./metafile.dl3n")

	io.Copy(w, metaF)
}

func StartServer() {
	http.HandleFunc("/upload", fileUploadHandler)
	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe("0.0.0.0:11112", nil)
}
