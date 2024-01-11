FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

RUN apk add --no-cache nodejs npm

RUN npm i

RUN npm run build

FROM scratch 

WORKDIR /app/

COPY --from=builder /app/main .

COPY --from=builder /app/dist/ ./dist/

CMD ["./main"]
