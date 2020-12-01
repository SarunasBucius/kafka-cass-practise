version = $$(git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X 'main.version=$$commit_hash' -extldflags -static"
build:
	docker build -t kcp-builder:$(version) ./builder;
	docker run --env commit_hash=$(version) --rm \
		--mount src=dep,dst=/go/pkg/mod/cache \
		--mount src=build-cache,dst=/root/.cache/go-build \
		--mount type=bind,src=$$PWD,dst=/usr/src/kcp \
		-w /usr/src/kcp kcp-builder:$(version) make build-binary;
	docker build -t kafka-cass-practise:$(version) ./bin;
	
build-binary:
	go build -tags "musl" $(LDFLAGS) -o ./bin/kcp ./cmd/kcp

run:
	docker run --rm -it kafka-cass-practise:$(version) 

start:
	env COMMIT_HASH=$(version) docker-compose up

stop:
	docker-compose down

rebuild:
	make stop;
	make build;
	make start;