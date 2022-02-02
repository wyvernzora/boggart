FROM alpine:3.15
RUN mkdir -p /etc/boggart/
WORKDIR /etc/boggart
COPY config.example.yml /etc/boggart/config.yml
COPY bin/boggart /bin/boggart
VOLUME ["/etc/boggart"]
CMD ["/bin/boggart"]
