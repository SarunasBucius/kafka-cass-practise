FROM golang:1.15.5-buster

# RUN apk add --no-cache librdkafka-dev    
CMD go build -o kcp -ldflags "-X 'main.version=$commit_hash' -extldflags -static" .
