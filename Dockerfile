FROM golang:1.21.5-alpine3.19

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]
