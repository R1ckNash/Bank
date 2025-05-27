FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/auth/go.mod services/auth/go.sum ./
COPY pkg ./pkg
COPY services/auth/internal ./internal
COPY services/auth/cmd ./cmd

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/auth

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service

ENV PORT=8080

CMD ["/app/auth-service"]