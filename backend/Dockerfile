FROM golang:1.20-alpine

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev

# Create logs directory
RUN mkdir -p /app/logs

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server

# Expose port
EXPOSE 8080

# Declare volume
VOLUME ["/app/logs"]

# Run the application
CMD ["./server"]