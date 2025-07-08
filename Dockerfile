FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

FROM debian:bookworm-slim

WORKDIR /root/
COPY --from=builder /app/app .

CMD ["./app"]
