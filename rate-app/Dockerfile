# syntax=docker/dockerfile:1

FROM golang:1.18-alpine 

WORKDIR /btc-app/rate-app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go get rate-api

COPY . .

RUN go build -o ./app ./cmd/rate-api \
    && chmod +x ./app

#EXPOSE 8080

ENTRYPOINT [ "./app" ]

