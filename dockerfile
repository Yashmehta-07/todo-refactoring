# FROM golang

# WORKDIR /app

# COPY go.mod go.sum ./

# COPY . .

# CMD [ "go", "run", "cmd/main.go" ]


# Stage 1: Build the Go application
FROM golang AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./

# Copy the application code
COPY . .

# Build the application (static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/app cmd/main.go


# Stage 2: Create a lightweight runtime image
FROM alpine:latest

WORKDIR /app

# Copy the statically built binary from the builder stage
COPY --from=builder /app/bin/app .

# Copy the migrations directory
COPY database/migrations /app/database/migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]


