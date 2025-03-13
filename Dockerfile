# Use official Golang image as a builder
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the application source code
COPY . .

# Set target architecture
ARG TARGETARCH

# Build the Go binary for the target architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags "-s -w -X main.version=$(git rev-parse --short HEAD)" -trimpath -o bin/bedrock-sidecar

# Create a minimal runtime image
FROM alpine:3.21.3

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/bin/bedrock-sidecar ./bin/bedrock-sidecar

# Run the binary
CMD ["./bin/bedrock-sidecar"]
