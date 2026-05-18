FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /car-service ./cmd/car-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /car-service .
COPY --from=builder /app/config ./config

EXPOSE 9998

ENTRYPOINT ["./car-service", "-config", "./config/docker.yml"]
