package gui

const index = `
<html>
    <body>
        <form
            enctype="multipart/form-data"
            action="http://localhost:8080/upload"
            method="post"
        >
            <input type="file" name="myFile" />
            <input type="submit" value="upload" />
        </form>
    </body>
</html>
`
