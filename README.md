# distributed-chaos-50.041

## Panopticon
The Panopticon is a debugger tool for viewing all chord nodes' predecessor and successor. Panopticon exists as a docker container that will sustain a consistent and easy to read table. Below is such an example 
 ```
                                     ID             predecessor               successor stabilized
panopticon    |           alpha (68601):     nodeCharlie (53816)       nodeBravo (96447)      true
panopticon    |       nodeBravo (96447):           alpha (68601)     nodeCharlie (53816)      true
panopticon    |     nodeCharlie (53816):       nodeBravo (96447)           alpha (68601)      true
```

To use this, you can run the following command in one terminal (append more services at the end if there are more). Or add a `-d` if you wish to detach the stdout.
```
docker-compose -f docker-compose.debug.yaml up alpha bravo charlie
```
Then in a clean terminal (which you don't need the previous print-outs anymore), run
```
docker-compose -f docker-compose.debug.yaml up panopticon
```
The above command will start the Panopticon container.  

## Protoc

### Installation
1. Install `protoc`: https://grpc.io/docs/protoc-installation/
2. Install go plugins for protocol compiler
    ```bash
    go get google.golang.org/protobuf/cmd/protoc-gen-go \
           google.golang.org/grpc/cmd/protoc-gen-go-grpc
    ```
3. Ensure that `$GOPATH/bin` is located within `PATH`

### Build/Rebuild gRPC library
Build the generated protoc code with:
```bash
protoc --go_out=./node/gossip/ --go_opt=paths=source_relative \
       --go-grpc_out=./node/gossip/ --go-grpc_opt=paths=source_relative \
       ./proto/internal.proto
```

<!--- 
for archiving
```bash
 protoc --go_out=plugins=gossip:../basic --go_opt=paths=source_relative basic.proto 
```
-->

## Docker

### Quick start (Build and Compose)
```bash
./builder.sh
s
```

### Build using DockerFile
To build a DockerFile as the image `chord_node`:
```bash
docker build -t chord_node .
```

Ensure that 4th step has been configured properly:
```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o node_exec <file with main function>
```

### docker-compose
```bash
docker-compose up
```
