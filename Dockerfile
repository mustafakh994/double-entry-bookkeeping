# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY migrations ./migrations
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 8080
CMD ["/app/main"]
