# Stage  1: Build the Go application
FROM golang:1.21.0

# Create a working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy the config.json file
COPY config.json /app/config.json


# Expose port 8080
EXPOSE 8080

# Build the application
RUN CGO_ENABLED=0 go build -o cmd/main cmd/main.go

# Start the application
CMD ["cmd/main"]
