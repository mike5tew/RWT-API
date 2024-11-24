# Stage 1: Build the Go application
FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main .
RUN ls -la /app  # Verify the main executable is created

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
RUN ls -la /root/  # Verify the main executable is copied
RUN chmod +x /root/main  # Ensure the main executable is executable

CMD ["./main"]