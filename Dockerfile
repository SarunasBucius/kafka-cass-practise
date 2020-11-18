ARG commit_hash
FROM golang
ARG commit_hash
ADD . /go/src/github.com/SarunasBucius/kafka-cass-practise
WORKDIR /go/src/github.com/SarunasBucius/kafka-cass-practise
RUN go get -d ./...
RUN go install -ldflags "-X 'main.version=${commit_hash}'" .

FROM ubuntu
COPY --from=0 /go/bin/kafka-cass-practise .
CMD ./kafka-cass-practise