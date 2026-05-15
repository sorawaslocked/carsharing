FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/user-service/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /user-service .
COPY config/docker.yml config/docker.yml

ENV CONFIG_PATH=/app/config/docker.yml

EXPOSE 9999

ENTRYPOINT ["/app/user-service"]
