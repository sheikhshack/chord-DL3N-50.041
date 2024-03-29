FROM golang as builder
# stagedbuild

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o node_exec ./cmd/node/node.go

FROM alpine
WORKDIR /built-app
COPY --from=builder /app/node_exec /built-app/
# ENVARS
ENV PEER_HOSTNAME=godbless
# defaults to 3 replicas!
ENV SUCCESSOR_LIST_SIZE=3
ENV LOG=info
ENV MY_PEER_DNS=DEFAULT
EXPOSE 9000 8888

CMD  ["./node_exec"]