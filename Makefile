version = $$(git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X 'main.version=$$commit_hash' -extldflags -static"
build:
	docker run --env commit_hash=$(version) --rm \
		--mount type=bind,src=$$GOPATH/pkg/mod/cache/download,dst=/go/pkg/mod/cache/download \
		--mount type=bind,src=$$PWD,dst=/usr/src/kcp \
		-w /usr/src/kcp golang:1.15.5-buster \
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