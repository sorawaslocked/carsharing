FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY protos/go.mod protos/go.sum ./protos/
COPY telemetry-service/go.mod telemetry-service/go.sum ./telemetry-service/

WORKDIR /app/telemetry-service
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY protos/ /app/protos/
COPY telemetry-service/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /telemetry ./cmd/server/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /telemetry .
COPY --from=builder /app/telemetry-service/docker.yml .

EXPOSE 50051

ENTRYPOINT ["./telemetry", "-config", "./docker.yml"]
