FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY internal internal

ENV CGO_ENABLED=0
RUN go build -o tg-gpt ./cmd/bot/main.go

FROM alpine:3.19.0

RUN apk update

COPY configs configs

COPY --from=builder /app/tg-gpt tg-gpt

CMD ["./tg-gpt"]
