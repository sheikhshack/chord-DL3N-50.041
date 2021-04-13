package gui

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// temporarily using a string for dev purposes
// will switch to using a static file later

const indexHtml = `
<html>
    <body>
        <form>
            <input type="file" id="fileUpload"/>
            <input type="submit"/>
        </form>

        <script>
            let file = null;
            const fileUpload = document.getElementById("fileUpload");
            fileUpload.addEventListener("change", () => {
                file = fileUpload.files[0];
            })
        </script>
    </body>
</html>
`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexHtml)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	inFile, inHandler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inFile.Close()

	fPath := "./" + inHandler.Filename
	f, err := os.Create(fPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, inFile)

	fmt.Fprintf(w, "file %s uploaded", inHandler.Filename)
}

func startServer() {
	http.HandleFunc("/upload", fileUploadHandler)
	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe("0.0.0.0:11112", nil)
}
