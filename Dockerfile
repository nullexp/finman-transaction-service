# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go application source code to the working directory
COPY . .

# Download and install any required Go modules
RUN go mod tidy

# Build the Go application
RUN go build -o finman-transaction-service ./cmd/

# Expose port 8082 to the outside world
EXPOSE 8082

# Run the executable
CMD ["./finman-transaction-service"]