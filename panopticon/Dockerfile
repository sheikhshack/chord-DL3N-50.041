FROM golang
# stagedbuild

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o node_exec ./panopticon/*.go

# using alpine for this, remember to use APK
#FROM alpine
#WORKDIR /built-app
#COPY --from=builder /app/node_exec /built-app/
EXPOSE 9000

CMD  ["./node_exec"]
