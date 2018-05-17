# Build
FROM golang:alpine AS build

RUN apk add --update git make build-base && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/khlieng/dispatch
COPY . /go/src/github.com/khlieng/dispatch
RUN go build .

# Runtime
FROM alpine

RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=build /go/src/github.com/khlieng/dispatch/dispatch /dispatch

EXPOSE 80/tcp

VOLUME ["/data"]

ENTRYPOINT ["/dispatch"]
CMD ["--dir=/data"]
