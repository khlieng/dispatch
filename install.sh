#!/usr/bin/env bash

Import="github.com/khlieng/dispatch/version"

Tag=$(git describe --tags)
Commit=$(git rev-parse --short HEAD)
Date=$(date +'%Y-%m-%dT%TZ')

go install -ldflags "-s -w -X $Import.Tag=$Tag -X $Import.Commit=$Commit -X $Import.Date=$Date"
