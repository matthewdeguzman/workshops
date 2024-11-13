# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

# Stage 2: Run
FROM golang:1.23

COPY --from=builder /app/bin/app /usr/local/bin/qrcode-backend

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/qrcode-backend"]
