FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/app /app

CMD ["/app"]
