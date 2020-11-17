version = $$(git rev-parse --short HEAD)
build:
	docker build -t kafka-cass-practise:$(version) .

run:
	docker run --rm --env commit_hash=$(version) kafka-cass-practise

start:
	docker-compose up --build

stop:
	docker-compose down