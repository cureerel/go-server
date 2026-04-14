# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Install build deps
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/build/cserver \
    ./cmd/server/main.go

# Stage 2: Minimal runtime
FROM scratch

# Carry timezone data and root certs from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/build/cserver /cserver

# Non-root user (uid 1001)
USER 1001

EXPOSE 8080

ENTRYPOINT ["/cserver"]
