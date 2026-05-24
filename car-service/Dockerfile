FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY protos/go.mod protos/go.sum ./protos/
COPY car-service/go.mod car-service/go.sum ./car-service/

WORKDIR /app/car-service
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY shared/ /app/shared/
COPY protos/ /app/protos/
COPY car-service/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /car-service ./cmd/car-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /car-service .
COPY --from=builder /app/car-service/config ./config

EXPOSE 9998

ENTRYPOINT ["./car-service", "-config", "./config/docker.yml"]
