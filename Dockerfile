FROM golang

ADD . /go/src/github.com/SarunasBucius/kafka-cass-practise

RUN go install github.com/SarunasBucius/kafka-cass-practise

ENTRYPOINT /go/bin/kafka-cass-practise
