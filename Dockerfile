FROM golang:1.20 AS builder
WORKDIR /app
copy . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]