# Use specific version of golang 
FROM golang:1.22.0-alpine AS builder

WORKDIR /app

ARG REACT_APP_API_URL=/api 

# Install git and gcc with specific alpine packages
RUN apk add --no-cache git gcc musl-dev

# Copy and download dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build with explicit CGO settings
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

# Use a specific alpine version for the final stage
FROM alpine:3.18

WORKDIR /root/
RUN mkdir -p /root/temp-images/
RUN apk add --no-cache bash

# Create necessary directories
RUN mkdir -p /app/fonts

# Copy fonts to the container
COPY ./fonts /app/fonts

COPY --from=builder /app/main .
COPY ./wait-for-it.sh /root/wait-for-it.sh
RUN chmod +x /root/main /root/wait-for-it.sh

CMD ["./wait-for-it.sh", "db:3306", "--", "./main", "-dbuser=${MYSQL_USER}", "-dbpassword=${MYSQL_PASSWORD}", "-dbname=${MYSQL_DATABASE}", "-dbhost=${MYSQL_HOST}"]


