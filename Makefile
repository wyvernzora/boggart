.PHONY: run docker clean

DOCKERFLAGS=--rm -it
PORT=2222

docker: bin/boggart
	docker build . -t boggart:dev

run: docker
	docker run $(DOCKERFLAGS) -v /etc/boggart:/etc/boggart -p $(PORT):2222 wyvernzora/boggart:dev

clean:
	rm -rf bin/*
	docker images prune

bin/boggart: *.go
	GOOS=linux go build -o bin/boggart
