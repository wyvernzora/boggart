FROM alpine:3.15

RUN mkdir -p /boggart/etc
WORKDIR /boggart

COPY config.example.yml etc/config.yml
COPY entrypoint.sh .
COPY boggart .

VOLUME ["/boggart/etc"]
CMD ["sh /boggart/entrypoint.sh"]
