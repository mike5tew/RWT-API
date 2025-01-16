FROM golang:1.23-alpine AS builder

WORKDIR /app

ARG REACT_APP_API_URL=/api 

# Install git and gcc
RUN apk add --no-cache git gcc musl-dev

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /app/main .

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

WORKDIR /root/
# Create necessary directories
RUN mkdir -p /root/temp-images/
RUN apk add --no-cache bash

# Create necessary directories
RUN mkdir -p /app/fonts

# Copy fonts to the container
COPY ./fonts /app/fonts

COPY --from=builder /app/main .
COPY ./wait-for-it.sh /root/wait-for-it.sh
RUN chmod +x /root/main /root/wait-for-it.sh

# Use wait-for-it script to wait for the database to be ready
CMD ["./wait-for-it.sh", "db:3306", "--", "./main", "-dbuser=${MYSQL_USER}", "-dbpassword=${MYSQL_PASSWORD}", "-dbname=${MYSQL_DATABASE}", "-dbhost=${MYSQL_HOST}"]


