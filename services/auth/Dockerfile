FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/auth

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/auth

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/auth-service /app/auth-service

CMD ["/app/auth-service"]