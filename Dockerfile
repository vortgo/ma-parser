FROM golang:1.12.0-alpine3.9

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --update --no-cache git

RUN go mod download

ADD . ./


RUN go build -o main .



CMD ["/app/main"]