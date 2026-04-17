# ===== Builder Stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy source
COPY src/ ./

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/api


# ===== Runtime Stage =====
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy timezone & certs (penting untuk production)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/app .

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]