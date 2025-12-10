FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rms ./cmd/server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/rms /usr/local/bin/rms
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/rms"]
