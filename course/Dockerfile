# Builder stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o main

# Final stage
FROM scratch

# Copy the executable binary from the builder stage into the final image
COPY --from=builder /app/main /main
COPY --from=builder /app/.env ./

# Run the application
CMD ["/main"]