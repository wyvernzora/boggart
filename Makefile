
image:
	docker build -t boggart:dev .

run: image
	docker run --rm -t -p 2222:2222 -v /etc/boggart:/etc/boggart boggart:dev
