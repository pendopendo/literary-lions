# Use the Go image to build the Go application
FROM golang:1.23.1 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o web ./cmd/web

# Start a new stage from scratch
FROM ubuntu:22.04

# Set Environment Variables
ENV PORT 8080

# Install SQLite
RUN apt-get update && apt-get install -y sqlite3

# Create and set a working directory inside the container
WORKDIR /app

# Copy the pre-built binary file, scripts and other required files from the previous stage
COPY --from=builder /app/web .
COPY --from=builder /app/scripts/initial-database.sql ./scripts/initial-database.sql
COPY --from=builder /app/ui/html ./ui/html
COPY --from=builder /app/ui/static ./ui/static

# Create a new SQLite database and populate it with initial data
RUN sqlite3 literarylionforum.db < ./scripts/initial-database.sql

# Command to run the executable
CMD ["./web"]