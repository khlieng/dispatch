#!/bin/sh
cp /etc/ssl/certs/ca-certificates.crt .
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/name_pending
docker build -t khlieng/name_pending .