FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY protos/go.mod protos/go.sum ./protos/
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/

WORKDIR /app/api-gateway
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY shared/ /app/shared/
COPY protos/ /app/protos/
COPY api-gateway/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/api-gateway

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /api-gateway ./api-gateway
COPY --from=builder /app/api-gateway/config ./config

EXPOSE 4000

CMD ["./api-gateway", "-config", "config/docker.yml"]
