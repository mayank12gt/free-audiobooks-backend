# Use the official golang image to create a build artifact.
FROM golang:1.21 AS builder

# Set the current working directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Use a lightweight base image.
FROM alpine:latest

# Set the working directory.
WORKDIR /root/

# Copy the pre-built binary from the previous stage.
COPY --from=builder /app/app .

# Copy .env file if needed.
COPY .env .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["./app"]
