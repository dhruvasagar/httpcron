.PHONY: default build run clean deploy

default: build

test:
	go test ./...

build: test
	docker build -t httpcron .

run: build
	docker run --rm -d -p "9000:9000" -it httpcron

deploy: build
	ssh -tt blog sudo service httpcron stop
	scp httpcron-linux blog:~/dsis.me/httpcron
	ssh -tt blog sudo service httpcron start
	rm -f httpcron-linux
