version = $$(git rev-parse --short HEAD)
build:
	docker build --build-arg commit_hash=$(version) -t kafka-cass-practise:$(version) .

run:
	docker run --rm kafka-cass-practise