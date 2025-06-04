FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/account

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/account-service ./cmd/account-service

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/account-service /app/account-service
COPY services/account/config/local.yaml ./config/local.yaml

CMD ["/app/account-service"]