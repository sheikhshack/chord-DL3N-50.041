FROM golang AS builder
# stagedbuild

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o node_exec ./cmd/node/node_join.go

# using alpine for this, remember to use APK
FROM alpine
WORKDIR /built-app
COPY --from=builder /app/node_exec /built-app/
# ENVARS
ENV NODE_ID=123
ENV PEER_HOSTNAME=godbless
EXPOSE 9000 8888

CMD  ["./node_exec"]