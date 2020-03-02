.PHONY: default test clean build run logs deploy

default: build

test:
	go test ./...

clean:
	docker stop httpcron &>/dev/null || true

build: clean test
	docker build -t httpcron .

run: build
	docker run \
		-d \
		-p "9000:9000" \
		-v 'dbdata:/dbdata' \
		--name httpcron \
		--restart always \
		-it httpcron

logs:
	docker logs -f httpcron

deploy: build
	ssh -tt blog sudo service httpcron stop
	scp httpcron-linux blog:~/dsis.me/httpcron
	ssh -tt blog sudo service httpcron start
	rm -f httpcron-linux
