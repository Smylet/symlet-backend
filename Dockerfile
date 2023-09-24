# Build stage
FROM golang:1.20-alpine3.16 AS builder

WORKDIR /app
COPY . .

# This step ensures that the directory structure under resources/env is preserved
RUN mkdir -p resources/env && mv app_test.env app.env resources/env/

RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .

# Copy the entire project directory to maintain the directory structure
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run
ENTRYPOINT ["/app/main"]
