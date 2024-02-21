FROM golang:1.22-alpine AS builder

COPY . /github.com/gogapopp/url-shortener/shortener/source/
WORKDIR /github.com/gogapopp/url-shortener/shortener/source/

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./bin/shortener cmd/shortener/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/gogapopp/url-shortener/shortener/source/bin/shortener .

CMD ["./shortener"]