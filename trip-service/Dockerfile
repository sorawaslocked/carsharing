FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY trip-service/go.mod trip-service/go.sum ./trip-service/

WORKDIR /app/trip-service
RUN go mod download

COPY shared/ /app/shared/
COPY trip-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /trip-service ./cmd/trip-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /trip-service .
COPY --from=builder /app/trip-service/config ./config

EXPOSE 9996

ENTRYPOINT ["./trip-service", "-config", "./config/docker.yml"]
