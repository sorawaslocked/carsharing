FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY protos/go.mod protos/go.sum ./protos/
COPY user-service/go.mod user-service/go.sum ./user-service/

WORKDIR /app/user-service
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY shared/ /app/shared/
COPY protos/ /app/protos/
COPY user-service/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/user-service/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /user-service .
COPY --from=builder /app/user-service/config/docker.yml config/docker.yml

ENV CONFIG_PATH=/app/config/docker.yml

EXPOSE 9999

ENTRYPOINT ["/app/user-service"]
