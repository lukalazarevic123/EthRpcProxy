# Start from the latest golang base image
FROM golang:1.21-alpine

# Install git
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download and install the dependencies
RUN go mod download

# Build the Go application
RUN go build -o main ./cmd

# Expose ports to the outside world
EXPOSE ${PORT}

# Command to run the executable
CMD ["./main"]
