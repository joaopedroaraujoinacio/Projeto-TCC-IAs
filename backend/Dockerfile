FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl openssl tzdata \
    && adduser -D -s /bin/sh appuser \
    && mkdir -p /app/certs \
    && openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /app/certs/server.key \
        -out /app/certs/server.crt \
        -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"
WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/ssl && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /app/ssl/server.key \
        -out /app/ssl/server.crt \
        -subj "/C=US/ST=State/L=City/O=MyOrg/OU=IT/CN=localhost" && \
    chmod 600 /app/ssl/server.key && \
    chmod 644 /app/ssl/server.crt && \
    chown -R appuser:appuser /app && \
    chmod +x /app/main
USER appuser

EXPOSE 8080

CMD ["./main"]

