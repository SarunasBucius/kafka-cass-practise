version = $$(git rev-parse --short HEAD)
build:
	docker build --build-arg commit_hash=$(version) -t kafka-cass-practise:$(version) . 

run:
	docker run --rm -it kafka-cass-practise:$(version) 

start:
	env COMMIT_HASH=$(version) docker-compose up
stop:
	docker-compose down