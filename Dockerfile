FROM golang
ARG commit_hash
WORKDIR /go/src/github.com/SarunasBucius/kafka-cass-practise
ADD . .
RUN go build -o kcp -ldflags "-X 'main.version=${commit_hash}' -linkmode external -extldflags -static" .

FROM scratch
COPY --from=0 /go/src/github.com/SarunasBucius/kafka-cass-practise/kcp .
CMD ["./kcp"]