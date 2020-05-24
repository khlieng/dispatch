#!/bin/sh -

Import="github.com/khlieng/dispatch/version"

Tag=$(git describe --tags)
Commit=$(git rev-parse --short HEAD)
Date=$(date +'%Y-%m-%dT%TZ')

CGO_ENABLED=0 go install -ldflags "-s -w -X $Import.Tag=$Tag -X $Import.Commit=$Commit -X $Import.Date=$Date"
