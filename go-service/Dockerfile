# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files (if using modules) and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code and the mapping file into the container
COPY . .

# Build the binary for Linux. Adjust the output binary name if needed.
RUN CGO_ENABLED=0 GOOS=linux go build -o go-service main.go

# Stage 2: Create a minimal container with the compiled binary
FROM alpine:latest

# Set the working directory in the minimal container
WORKDIR /root/

# Copy the binary and the mapping file from the builder stage to the current stage.
COPY --from=builder /app/go-service .
COPY --from=builder /app/mapping.json .

# Expose port 8080 for the service
EXPOSE 8080

# Command to run the binary when the container starts
CMD ["./go-service"]
