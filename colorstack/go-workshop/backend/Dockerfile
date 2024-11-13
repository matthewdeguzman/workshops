# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN make build
EXPOSE 8080

ENTRYPOINT ["/app/bin/app"]
