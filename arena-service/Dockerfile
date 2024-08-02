FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/arena-backend

RUN CGO_ENABLED=0 GOOS=linux go build -o /arena .

EXPOSE 8000

CMD ["/arena"]
