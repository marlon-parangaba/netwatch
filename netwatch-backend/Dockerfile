# --- Build Stage ---
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o netwatch-server ./cmd/server

# --- Run Stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/netwatch-server .
COPY --from=builder /app/config.yaml .

# Migrations are embedded into the binary via go:embed
COPY --from=builder /app/migrations ./migrations

RUN addgroup -S netwatch && adduser -S netwatch -G netwatch
RUN mkdir -p /app/logs && chown -R netwatch:netwatch /app
USER netwatch

EXPOSE 8080
EXPOSE 162/udp

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./netwatch-server"]
