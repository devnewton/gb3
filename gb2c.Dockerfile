FROM golang:1.16 AS gb2c_builder
WORKDIR /go/src/github.com/devnewton/gb3
COPY . .
RUN go build ./cmd/gb2c

FROM debian:stable as gb2c
WORKDIR /root/
RUN apt update && apt install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
COPY --from=gb2c_builder /go/src/github.com/devnewton/gb3/gb2c ./
EXPOSE 16666
CMD ["./gb2c"]