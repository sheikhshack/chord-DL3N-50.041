# distributed-chaos-50.041

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
