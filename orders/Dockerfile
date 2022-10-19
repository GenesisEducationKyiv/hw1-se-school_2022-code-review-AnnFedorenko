FROM golang:1.18-alpine as builder

WORKDIR /go/src/orders

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/orders \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/orders/app .
COPY --from=builder /go/src/orders/.env .

ENTRYPOINT [ "./app" ]
