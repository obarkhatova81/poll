FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
EXPOSE 8081

CMD ["./main"]
