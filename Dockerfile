# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy Source Code
COPY . .

# Build API
RUN go build -o pesu-api ./api/main.go

# Build Worker
RUN go build -o pesu-worker ./worker/main.go

# Runtime Stage (API)
FROM alpine:latest AS api
WORKDIR /root/
COPY --from=builder /app/pesu-api .
EXPOSE 8080
CMD ["./pesu-api"]

# Runtime Stage (Worker)
FROM alpine:latest AS worker
WORKDIR /root/
# Install Docker CLI (needed to spawn containers)
RUN apk add --no-cache docker-cli
COPY --from=builder /app/pesu-worker .
CMD ["./pesu-worker"]
