# Build stage
FROM golang:1.23.1-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init -g main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Define build arguments
ARG PORT
ARG DB_HOST
ARG DB_PORT
ARG DB_DATABASE
ARG DB_USERNAME
ARG DB_PASSWORD
ARG DB_ROOT_PASSWORD
ARG REDIS_HOST
ARG REDIS_PORT
ARG JWT_SECRET

# Set environment variables from build arguments
ENV PORT=${PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_DATABASE=${DB_DATABASE}
ENV DB_USERNAME=${DB_USERNAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
ENV REDIS_HOST=${REDIS_HOST}
ENV REDIS_PORT=${REDIS_PORT}
ENV JWT_SECRET=${JWT_SECRET}

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE ${PORT}

# Run the application
CMD ["./main"] 