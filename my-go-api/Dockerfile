# Dockerfile

# Use official Go image
FROM golang:1.23

# Set working directory inside container
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod tidy

# Copy all files from current dir to container
COPY . .

# Download dependencies
RUN go mod download

# Set default command (can be overridden in docker-compose)
CMD ["go", "run", "main.go", "db.go"]
