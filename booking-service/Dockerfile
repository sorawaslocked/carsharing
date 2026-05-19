FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY booking-service/go.mod booking-service/go.sum ./booking-service/

WORKDIR /app/booking-service
RUN go mod download

COPY shared/ /app/shared/
COPY booking-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /booking-service ./cmd/booking-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /booking-service .
COPY --from=builder /app/booking-service/config ./config

EXPOSE 9997

ENTRYPOINT ["./booking-service", "-config", "./config/docker.yml"]
