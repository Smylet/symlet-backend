# Build stage
FROM golang:1.20-alpine3.16 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .

COPY . .
COPY app.env app.env
# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run
ENTRYPOINT ["/app/main"]