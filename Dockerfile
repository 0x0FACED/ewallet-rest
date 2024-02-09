FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

RUN mkdir config

COPY --from=builder /app/main .

COPY --from=builder /app/config/server.toml config/

CMD ["./main"]
