version = $$(git rev-parse --short HEAD)
build:
	docker run --env commit_hash=$(version) --rm -v "$$PWD":/usr/src/kcp -w /usr/src/kcp golang make build-binary;
	docker build -t kafka-cass-practise:$(version) .
	
build-binary:
	CGO_ENABLED=0 go build -o kcp -ldflags "-X 'main.version=$$commit_hash'" .

run:
	docker run --rm -it kafka-cass-practise:$(version) 

start:
	env COMMIT_HASH=$(version) docker-compose up
stop:
	docker-compose down