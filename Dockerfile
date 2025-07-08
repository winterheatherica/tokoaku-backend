FROM gcr.io/distroless/base

WORKDIR /

COPY app-railway /app

CMD ["/app"]
