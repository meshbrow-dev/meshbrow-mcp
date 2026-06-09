FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o meshbrow-mcp ./cmd/mcp-server

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/meshbrow-mcp /usr/local/bin/meshbrow-mcp
ENTRYPOINT ["meshbrow-mcp"]
