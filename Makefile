.PHONY: run docker

docker: bin/boggart
	docker build . -t boggart:dev

bin/boggart: *.go
	GOOS=linux go build -o bin/boggart
