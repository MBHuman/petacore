FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Компилируем API сервер
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server cmd/api/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-server .

EXPOSE 8080

CMD ["./api-server"]
