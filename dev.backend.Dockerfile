FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY backend/go.mod ./

RUN go mod download

COPY backend/ ./

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
