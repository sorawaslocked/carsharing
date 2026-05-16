FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /booking-service ./cmd/booking-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /booking-service .
COPY --from=builder /app/config ./config

EXPOSE 9997

ENTRYPOINT ["./booking-service", "-config", "./config/docker.yml"]
