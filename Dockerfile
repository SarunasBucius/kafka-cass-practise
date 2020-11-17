FROM golang
ADD . /go/src/github.com/SarunasBucius/kafka-cass-practise
RUN go get -d ./...
RUN go install github.com/SarunasBucius/kafka-cass-practise

FROM ubuntu
COPY --from=0 /go/bin/kafka-cass-practise .
CMD ./kafka-cass-practise -ldflags="-X 'main.version=$(commit_hash)'"