FROM golang:1.16 AS gb0_builder
WORKDIR /go/src/github.com/devnewton/gb3
COPY . .
RUN go build ./cmd/gb0

FROM debian:stable as gb0
WORKDIR /root/
RUN apt update && apt install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
COPY --from=gb0_builder /go/src/github.com/devnewton/gb3/gb0 ./
EXPOSE 16667
CMD ["./gb0"]