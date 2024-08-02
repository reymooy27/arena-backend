FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/auth-service

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .

EXPOSE 8000

CMD ["/auth-service"]
