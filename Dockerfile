# Используйте официальный образ Golang для сборки
FROM golang:latest as builder

# Установите рабочую директорию внутри контейнера
WORKDIR /app

# Скопируйте go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Загрузите все зависимости
RUN go mod download

# Скопируйте исходный код в рабочую директорию
COPY . .

# Соберите приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Используйте образ alpine для выполнения
FROM alpine:latest

# Установите рабочую директорию внутри контейнера
WORKDIR /root/

RUN mkdir config

COPY --from=builder /app/main .

COPY --from=builder /app/config/server.toml config/

# Запустите приложение
CMD ["./main"]
