FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service /app/cmd/auth-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth-service .
COPY --from=builder /app/db/migrations /root/db/migrations
# COPY --from=builder /app/.env .

EXPOSE 8001

CMD ["./auth-service"]
