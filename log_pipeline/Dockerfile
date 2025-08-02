# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main cmd/grpc/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose the gRPC port
EXPOSE 50051

# Expose the port your Go application might listen on (if it's a web server)
# EXPOSE 8080 # Uncomment if your Go app is a web server

#  Define environment variables
ENV POSTGRES_HOST=postgres_db
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=logpipeline

# Command to run the executable
CMD ["./main"]
