version = $$(git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X 'main.version=$$commit_hash' -extldflags -static"
build:
	# docker run --env commit_hash=$(version) --rm -v "$$PWD":/usr/src/kcp -w /usr/src/kcp go-builder;
	docker run --env commit_hash=$(version) --rm \
		-v "$$PWD":/usr/src/kcp -w /usr/src/kcp golang:1.15.5-buster \
		make build-binary;
	docker build -t kafka-cass-practise:$(version) .;
	rm -f ./kcp
	
build-binary:
	go build $(LDFLAGS) -o kcp .

run:
	docker run --rm -it kafka-cass-practise:$(version) 

start:
	env COMMIT_HASH=$(version) docker-compose up
stop:
	docker-compose down