# name_pending
Web-based IRC client in Go.

## Building

#### Requirements
* [Go](http://golang.org/doc/install)
* [Node.js + npm](http://nodejs.org/download/)

#### Get the source
```bash
go get github.com/khlieng/name_pending
```

#### Compile the server
```bash
cd $GOPATH/src/github.com/khlieng/name_pending
go build -o bin/name_pending
```

#### Build the client
```bash
npm install -g gulp

cd client
npm install
gulp -p
```
