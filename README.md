# name_pending
Web-based IRC client in Go.

## Installing
```bash
go get github.com/khlieng/name_pending
```

## Running
```bash
name_pending
```

## Building the server

#### Requirements
* [Go](http://golang.org/doc/install)

```bash
cd $GOPATH/src/github.com/khlieng/name_pending
go install
```

## Building the client

#### Requirements
* [Node.js + npm](https://nodejs.org/download/)

```bash
npm install -g gulp
go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...

cd $GOPATH/src/github.com/khlieng/name_pending/client
npm install
gulp -p
go-bindata-assetfs -nomemcopy dist/...
mv bindata_assetfs.go ../bindata.go

# Rebuild the server :)
```