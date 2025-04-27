# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build ./cmd/main.go

RUN mkdir -p /logs && chmod 755 /logs

# Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .

CMD [ "/app/main" ]