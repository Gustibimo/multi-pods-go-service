# Stage 1: Build the Go application
FROM golang:1.21.5-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git and other dependencies if needed
RUN #apk add --no-cache git

# Copy go mod and sum files
COPY ../go.mod go.sum ./

# Download dependencies (add verbose output for debugging)
RUN go mod download -x

# Copy the rest of the application code
COPY . .

ENV GOCACHE=/target/.cache/go-build
# Build the Go application
RUN  go build -o bom-import-backend
# Stage 2: Create the runtime container
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/bom-import-backend .
COPY bom.xlsx .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./bom-import-backend", "http-server"]
