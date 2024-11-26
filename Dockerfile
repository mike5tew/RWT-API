# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Install git and gcc
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /app/main .
RUN ls -la /app  # Verify the main executable is created

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache bash

COPY --from=builder /app/main .
COPY wait-for-it.sh /root/wait-for-it.sh
RUN chmod +x /root/main /root/wait-for-it.sh

CMD ["./wait-for-it.sh", "db:3306", "--timeout=60", "--", "./main"]