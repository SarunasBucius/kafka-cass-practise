ARG commit_hash
FROM golang
ADD . /go/src/github.com/SarunasBucius/kafka-cass-practise
RUN go install github.com/SarunasBucius/kafka-cass-practise

FROM scratch
ARG commit_hash
ENV commit_hash=$commit_hash
COPY --from=0 /go/bin/kafka-cass-practise .
CMD ["./kafka-cass-practise"]