# Use an official Golang runtime as the base image
FROM golang:1.17-alpine

# Install Git
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Download and cache Go modules
RUN go mod download

# Build the Go application
RUN go build -o producer ./src/main.go

# Expose any necessary ports
EXPOSE 8081

# Run the Go application
CMD ["./producer"]
