FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/api /app/api

EXPOSE 8888

CMD ["/app/api"]