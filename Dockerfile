FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/main .

RUN chown appuser:appuser /app/main
RUN chmod +x /app/main

USER appuser

EXPOSE 8080

CMD ["./main"]
