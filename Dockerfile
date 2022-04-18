FROM golang:1.18-alpine3.15 as builder

WORKDIR /go/src/github.com/wyvernzora/boggart
COPY . .

RUN apk --no-cache add git alpine-sdk
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o boggart cmd/boggart/*.go && chmod +x boggart
RUN mkdir -p /etc/boggart && cp configs/config.example.yml /etc/boggart/config.yml


FROM scratch

WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/wyvernzora/boggart/boggart /root/boggart
COPY --from=builder /etc/boggart /etc/boggart

VOLUME ["/etc/boggart"]
CMD ["/root/boggart"]
