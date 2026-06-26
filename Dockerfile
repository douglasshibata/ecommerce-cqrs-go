FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY internal/eventstore/schema.sql ./internal/eventstore/
COPY internal/projection/read_schema.sql ./internal/projection/

EXPOSE 8080

CMD ["./main"]
