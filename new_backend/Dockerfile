# Builder stage - Use Debian instead of Alpine
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install build dependencies (Debian uses apt instead of apk)
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copy source code
COPY . .

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Runtime stage - Use Debian slim instead of Alpine
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary and assets
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
