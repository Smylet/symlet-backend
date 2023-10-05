# Build stage
FROM golang:1.20-alpine3.16 AS builder

WORKDIR /app
COPY . .
ENV GOCACHE=/root/.cache/go-build

# This step ensures that the directory structure under resources/env is preserved
# RUN mkdir -p resources/env && mv app_test.env app.env resources/env/

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o main main.go

# Intermediate stage to copy migrations
FROM arigaio/atlas:latest AS atlas

COPY migrations /migrations

# Run stage
FROM alpine:3.16
WORKDIR /app

# Copy from the builder stage
COPY --from=builder /app/main .

# Copy migrations from atlas stage
COPY --from=atlas /migrations /migrations

# Copy the entire project directory to maintain the directory structure
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run
ENTRYPOINT ["/app/main"]
