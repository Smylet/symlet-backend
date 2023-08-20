# Stage 1: Build
FROM golang:1.20-alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy binary from build to main folder
COPY --from=builder /build/main /app/

# Command to run
ENTRYPOINT ["/app/main"]