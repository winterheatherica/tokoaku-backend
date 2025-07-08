FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app ./cmd

FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /app/app /app

CMD ["/app"]
