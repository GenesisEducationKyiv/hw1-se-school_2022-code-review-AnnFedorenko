FROM golang:1.18-alpine

WORKDIR /btc-app/log-consumer

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go get log-consumer

COPY . .

RUN go build -o ./app ./cmd/log-consumer \
    && chmod +x ./app

ENTRYPOINT [ "./app" ]