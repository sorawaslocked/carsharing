FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY protos/go.mod protos/go.sum ./protos/
COPY trip-service/go.mod trip-service/go.sum ./trip-service/

WORKDIR /app/trip-service
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY shared/ /app/shared/
COPY protos/ /app/protos/
COPY trip-service/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /trip-service ./cmd/trip-service/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /trip-service .
COPY --from=builder /app/trip-service/config ./config

EXPOSE 9996

ENTRYPOINT ["./trip-service", "-config", "./config/docker.yml"]
