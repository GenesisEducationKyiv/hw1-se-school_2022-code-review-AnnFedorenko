FROM golang:1.18-alpine as builder

WORKDIR /go/src/customers

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/customers \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/customers/app .
COPY --from=builder /go/src/customers/.env .

ENTRYPOINT [ "./app" ]
