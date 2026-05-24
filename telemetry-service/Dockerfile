FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /telemetry ./cmd/server/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /telemetry .
COPY --from=builder /app/docker.yml .

EXPOSE 50051

ENTRYPOINT ["./telemetry", "-config", "./docker.yml"]
