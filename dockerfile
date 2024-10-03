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

# Install SQLite and initialize the database
RUN apt-get update && apt-get install -y sqlite3
COPY ./scripts/initial-database.sql ./scripts/initial-database.sql
RUN sqlite3 /app/literarylionforum.db < ./scripts/initial-database.sql

# Start a new stage from scratch
FROM ubuntu:22.04

# Set Environment Variables
ENV PORT 8080

# Create and set a working directory inside the container
WORKDIR /app

# Declare a volume for the SQLite database
VOLUME /app/data

# Copy the pre-built binary file and initialized database from the builder stage
COPY --from=builder /app/web .
COPY --from=builder /app/literarylionforum.db /app/data/literarylionforum.db
COPY --from=builder /app/ui/html ./ui/html
COPY --from=builder /app/ui/static ./ui/static

# Command to run the executable
CMD ["./web"]