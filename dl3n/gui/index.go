package gui

const index = `
<!doctype html>
<html lang="en">

<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <title>DL3N Demo</title>

    <style>
        html,
        body {
            height: 100%;
        }

        .file-preview {
            width: 15%;
            background-color: lightgrey;
            background-image: url("data:image/svg+xml,%3C%3Fxml version='1.0' encoding='iso-8859-1'%3F%3E%3C!-- Generator: Adobe Illustrator 18.0.0, SVG Export Plug-In . SVG Version: 6.00 Build 0) --%3E%3C!DOCTYPE svg PUBLIC '-//W3C//DTD SVG 1.1//EN' 'http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd'%3E%3Csvg version='1.1' id='Capa_1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink' x='0px' y='0px' viewBox='0 0 317.001 317.001' style='enable-background:new 0 0 317.001 317.001;' xml:space='preserve'%3E%3Cpath d='M270.825,70.55L212.17,3.66C210.13,1.334,207.187,0,204.093,0H55.941C49.076,0,43.51,5.566,43.51,12.431V304.57 c0,6.866,5.566,12.431,12.431,12.431h205.118c6.866,0,12.432-5.566,12.432-12.432V77.633 C273.491,75.027,272.544,72.51,270.825,70.55z M55.941,305.073V12.432H199.94v63.601c0,3.431,2.78,6.216,6.216,6.216h54.903 l0.006,222.824H55.941z'/%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3Cg%3E%3C/g%3E%3C/svg%3E%0A");
            background-size: 30%;
            background-repeat: no-repeat;
            background-position: center;
            border-radius: 1rem;
        }


        .file-upload {
            flex-grow: 1;
            padding-top: 2.5rem;
            padding-bottom: 2.5rem;
            padding-left: 1rem;
            border-radius: 1rem;
            border: grey 0.2rem dashed;
            margin-left: 0.2rem;
        }
    </style>
</head>

<body class="d-flex flex-column align-items-center justify-content-center">
    <div class="container ">
        <h1 class="text-center">DL3N Demo</h1>
        <div class="card">
            <div class="card-header">
                <p class="mb-0">Current State: <b><a id="currentState">WAITING</a></b></p>
            </div>
            <div class="card-body">
                <ul class="nav nav-tabs">
                    <li class="nav-item">
                        <a href="#seed" class="nav-link active" data-toggle="tab">Seed File</a>
                    </li>
                    <li class="nav-item">
                        <a href="#get" class="nav-link" data-toggle="tab">Get File</a>
                    </li>
                </ul>
                <div class="tab-content container py-2">

                    <!-- Seed Tab Start -->
                    <div class="tab-pane show active" id="seed">
                        <h4>Upload file to seed</h4>

                        <form id="fileUpload" enctype="multipart/form-data" action="/upload" method="post">
                            <div class="d-flex flex-row my-2">
                                <img id="file-preview" src="" alt="" class="file-preview" />
                                <input name="uploadFile" type="file" class="file-upload">
                            </div>
                        </form>

                        <div class="d-flex flex-row justify-content-center row my-2">
                            <button id="generateMeta" class="btn btn-primary mx-1">Generate Metadata File</button>
                            <button id="downloadMeta" class="btn btn-primary mx-1" disabled>Download Metadata
                                File</button>
                        </div>

                        <div class="d-flex flex-row justify-content-center row my-2">
                            <button id="startSeed" class="btn btn-success mx-1" disabled> Start Seeding </button>
                            <button id="stopSeed" class="btn btn-danger mx-1" disabled>Stop Seeding</button>
                        </div>
                    </div>
                    <!-- Seed Tab End -->

                    <!-- Get Tab Start -->
                    <div class="tab-pane" id="get">
                        <h4>Upload file to seed</h4>

                        <form id="metaUpload" enctype="multipart/form-data" action="/upload" method="post">
                            <div class="d-flex flex-row my-2">
                                <img id="file-preview" src="" alt="" class="file-preview" />
                                <input name="uploadMeta" type="file" class="file-upload">
                            </div>
                        </form>

                        <div class="d-flex flex-row justify-content-center row my-2">
                            <button id="startGet" class="btn btn-primary mx-1" disabled>Start Download</button>
                            <button id="getFile" class="btn btn-success mx-1" disabled>Save File</button>
                        </div>
                    </div>
                    <!-- Get Tab End -->
                </div>
            </div>
        </div>
    </div>


    <!-- jQuery first, then Popper.js, then Bootstrap JS -->
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js"
        integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo"
        crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"
        integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1"
        crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"
        integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM"
        crossorigin="anonymous"></script>

    <script>
        // Change to empty in index.go
        const BASE_URL = ""

        const state = {
            State: "WAITING",
            SeedMeta: null,
            GetMeta: null,
        };

        const setState = (s) => {
            state.State = s.State;
            state.SeedMeta = s.SeedMeta;
            state.GetMeta = s.GetMeta;

            $("#downloadMeta").attr("disabled", !state.SeedMeta);
            $("#startSeed").attr("disabled", !state.SeedMeta || state.State == "SEEDING");
            $("#stopSeed").attr("disabled", state.State != "SEEDING");
            $("#currentState").html(state.State);

            $("#startGet").attr("disabled", !state.GetMeta || state.State == "GETTING");
            $("#getFile").attr("disabled", state.State != "GET_DONE");
        }

        const req = (endpoint, body) => {
            fetch(BASE_URL + endpoint, { method: "POST", body })
                .then(r => r.json())
                .then(s => { setState(s) })
        }

        const fileUploadForm = $("#fileUpload");

        $("#fileUpload").change((e) => {
            e.preventDefault();
            const formData = new FormData($("#fileUpload")[0]);
            const file = formData.get("uploadFile");
            const url = URL.createObjectURL(file);
            $("#file-preview").attr("src", url)
        })

        $("#generateMeta").click((e) => {
            e.preventDefault();
            const formData = new FormData($("#fileUpload")[0]);
            req("/upload", formData);
        })

        $("#downloadMeta").click((e) => {
            e.preventDefault();

            const text = JSON.stringify(state.SeedMeta);
            const element = document.createElement("a");

            element.setAttribute("href", "data:text/plain;charset=utf-8," + encodeURIComponent(text));
            element.setAttribute("download", state.SeedMeta.Name + ".dl3n");
            element.style.display = "none";
            document.body.appendChild(element);
            element.click();
            document.body.removeChild(element);
        })

        $("#startSeed").click((e) => {
            req("/startSeed", null);
        });

        $("#stopSeed").click((e) => {
            req("/stopSeed", null);
        });

        $("#metaUpload").change((e) => {
            e.preventDefault();
            const formData = new FormData($("#metaUpload")[0]);
            req("/uploadMeta", formData);
        })

        $("#startGet").click((e) => {
            req("/startGet", null);
        });

        $("#getFile").click((e) => {
            const element = document.createElement("a");
            element.setAttribute("href", BASE_URL + "/getFile");
            element.setAttribute("target", "_blank");
            element.setAttribute("download", state.GetMeta.Name);
            element.style.display = "none";
            document.body.appendChild(element);
            element.click();
            document.body.removeChild(element);
        });

        const pollState = () => {
            req("/getState", null);
            setTimeout(pollState, 1000);
        }
        pollState();
    </script>
</body>

</html>
`
