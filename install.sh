#!/usr/bin/env bash

Import="github.com/khlieng/dispatch/commands"

Version=$(git describe --tags)
Commit=$(git rev-parse --short HEAD)
Date=$(date +'%Y-%m-%dT%TZ')

go install -ldflags "-s -w -X $Import.version=$Version -X $Import.commit=$Commit -X $Import.date=$Date"
