FROM golang

ADD . /go/src/github.com/khlieng/name_pending

RUN go install github.com/khlieng/name_pending

VOLUME ["/data"]

ENTRYPOINT ["/go/bin/name_pending"]
CMD ["-p=8080", "--dir=/data"]