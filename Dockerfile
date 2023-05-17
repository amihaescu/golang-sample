# Stage 1: Build Stage
FROM golang:1.19-alpine3.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files to the working directory
COPY src/ .

# Fetch dependencies
RUN go mod download

# Build the Go application
RUN go build -o main

# Stage 2: Final Image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which the application will run
EXPOSE 8080

# Set the entry point for the container
CMD ["./main"]
