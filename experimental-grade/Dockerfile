FROM golang AS builder
# stagedbuild

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o mike_exec .

# using alpine for this, remember to use APK
FROM alpine
WORKDIR /built-app
COPY --from=builder /app/mike_exec /built-app/
# ENVARS
EXPOSE 1234

CMD  ["./mike_exec"]