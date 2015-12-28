#!/bin/sh
cp /etc/ssl/certs/ca-certificates.crt .
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags -w -o build/dispatch
docker build -t khlieng/dispatch .
