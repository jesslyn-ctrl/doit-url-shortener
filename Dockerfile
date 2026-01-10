# ====================
# 1. Build Stage
# ====================
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install git (might be needed for go modules)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd

# ====================
# 2. Runtime Stage
# ====================
FROM gcr.io/distroless/base-debian12

# Run as non-root user
USER nonroot:nonroot

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]