# Start from the official Go image
FROM golang:1.22
ENV GOTOOLCHAIN=auto

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o messenger ./cmd

# Expose the app port
EXPOSE 8080

# Command to run the executable
CMD ["./messenger"]
