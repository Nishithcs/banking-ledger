# Use the official Golang image as a base
FROM golang:1.23.8 AS builder

# Set working directory
WORKDIR /app

COPY go.mod go.sum ./

# Set GOPROXY to direct to avoid proxy timeouts
# ENV GOPROXY=off
ENV GOPROXY=https://proxy.golang.org,direct


RUN apt-get update && apt-get install -y git ca-certificates

RUN go mod download

# Copy the source code
COPY . .


# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker


# Use a smaller image for the final container
FROM alpine:latest

# Install CA certificates for HTTPS
# RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /. .

# Expose the port the API will run on
EXPOSE 8080

# # Command to run the executable
# CMD ["./main"]
CMD ["sh", "-c", "while true; do sleep 3600; done"]
