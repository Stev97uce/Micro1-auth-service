# Etapa de compilación
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

# Etapa de ejecución
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .env

EXPOSE 8080

CMD ["./main"]
