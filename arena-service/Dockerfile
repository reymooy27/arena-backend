FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o arena-service /app/cmd/arena-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/arena-service .
COPY --from=builder /app/db/migrations /root/db/migrations
# COPY --from=builder /app/.env .

EXPOSE 50002

CMD ["./arena-service"]
