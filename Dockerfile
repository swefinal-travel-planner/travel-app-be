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
ARG MAIL_HOST
ARG MAIL_PORT
ARG MAIL_USERNAME
ARG MAIL_PASSWORD
ARG MAIL_FROM
ARG MAIL_FROM_NAME
ARG REDIS_HOST
ARG REDIS_PORT
ARG REDIS_PASSWORD
ARG JWT_SECRET
ARG ALLOWED_ORIGIN

# Set environment variables from build arguments
ENV PORT=${PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_DATABASE=${DB_DATABASE}
ENV DB_USERNAME=${DB_USERNAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
ENV MAIL_HOST=${MAIL_HOST}
ENV MAIL_PORT=${MAIL_PORT}
ENV MAIL_USERNAME=${MAIL_USERNAME}
ENV MAIL_PASSWORD=${MAIL_PASSWORD}
ENV MAIL_FROM=${MAIL_FROM}
ENV MAIL_FROM_NAME=${MAIL_FROM_NAME}
ENV REDIS_HOST=${REDIS_HOST}
ENV REDIS_PORT=${REDIS_PORT}
ENV REDIS_PASSWORD=${REDIS_PASSWORD}
ENV JWT_SECRET=${JWT_SECRET}
ENV ALLOWED_ORIGIN=${ALLOWED_ORIGIN}

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE ${PORT}

# Run the application
CMD ["./main"] 