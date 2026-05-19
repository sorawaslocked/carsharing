FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/

WORKDIR /app/api-gateway
RUN go mod download

COPY shared/ /app/shared/
COPY api-gateway/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/api-gateway

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /api-gateway ./api-gateway
COPY --from=builder /app/api-gateway/config ./config

EXPOSE 4000

CMD ["./api-gateway", "-config", "config/docker.yml"]
