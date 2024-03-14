FROM golang:1.22 AS builder

RUN apt-get update && apt-get install -y gcc g++

COPY . /github.com/gogapopp/url-shortener/metrics/source/
WORKDIR /github.com/gogapopp/url-shortener/metrics/source/

RUN go mod download
RUN CGO_ENABLED=1 go build -o ./bin/metrics cmd/metrics/main.go

FROM debian:latest

WORKDIR /root/
COPY --from=builder /github.com/gogapopp/url-shortener/metrics/source/bin/metrics .

CMD ["./metrics"]