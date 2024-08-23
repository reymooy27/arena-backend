FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service /app/cmd/payment-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/payment-service .
COPY --from=builder /app/db/migrations /root/db/migrations

EXPOSE 50001

CMD ["./payment-service"]
