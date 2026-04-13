# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# CGO required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./
RUN go mod tidy
RUN go mod download

COPY . .

# Swagger docs are pre-generated
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o golearn .

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/golearn .

# SQLite data volume mount point
RUN mkdir -p /data

EXPOSE 8090

CMD ["./golearn"]
