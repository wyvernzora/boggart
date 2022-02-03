.PHONY: run docker clean

DOCKERFLAGS=--rm -it
PORT=2222

docker: bin/boggart
	docker build ./docker -t boggart:dev

run: docker
	docker run $(DOCKERFLAGS) -v /etc/boggart:/boggart/etc -p $(PORT):2222 wyvernzora/boggart:dev

clean:
	rm -rf docker/boggart
	docker images prune

bin/boggart: *.go
	GOOS=linux go build -o docker/boggart
