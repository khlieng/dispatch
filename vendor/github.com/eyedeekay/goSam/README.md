goSam
=====

A go library for using the [I2P](https://geti2p.net/en/) Simple Anonymous
Messaging ([SAM version 3.0](https://geti2p.net/en/docs/api/samv3)) bridge. It
has support for all streaming features SAM version 3.2.

This is widely used and easy to use, but thusfar, mostly by me. It sees a lot of
testing and no breaking changes to the API are expected.

## Installation
```
go get github.com/eyedeekay/goSam
```

## Using it for HTTP Transport

`Client.Dial` implements `net.Dial` so you can use go's library packages like http.

```go
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cryptix/goSam"
)

func main() {
	// create a default sam client
	sam, err := goSam.NewDefaultClient()
	checkErr(err)

	log.Println("Client Created")

	// create a transport that uses SAM to dial TCP Connections
	tr := &http.Transport{
		Dial: sam.Dial,
	}

	// create  a client using this transport
	client := &http.Client{Transport: tr}

	// send a get request
	resp, err := client.Get("http://stats.i2p/")
	checkErr(err)
	defer resp.Body.Close()

	log.Printf("Get returned %+v\n", resp)

	// create a file for the response
	file, err := os.Create("stats.html")
	checkErr(err)
	defer file.Close()

	// copy the response to the file
	_, err = io.Copy(file, resp.Body)
	checkErr(err)

	log.Println("Done.")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```

## Using it as a SOCKS proxy

`client` also implements a resolver compatible with
[`getlantern/go-socks5`](https://github.com/getlantern/go-socks5),
making it very easy to implement a SOCKS5 server.

```go
package main

import (
  "flag"

	"github.com/eyedeekay/goSam"
	"github.com/getlantern/go-socks5"
	"log"
)

var (
  samaddr = flag.String("sam", "127.0.0.1:7656", "SAM API address to use")
  socksaddr = flag.String("socks", "127.0.0.1:7675", "SOCKS address to use")
)

func main() {
	sam, err := goSam.NewClient(*samaddr)
	if err != nil {
		panic(err)
	}
	log.Println("Client Created")

	// create a transport that uses SAM to dial TCP Connections
	conf := &socks5.Config{
		Dial:     sam.DialContext,
		Resolver: sam,
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", *socksaddr); err != nil {
		panic(err)
	}
}
```

### .deb package

A package for installing this on Debian is buildable, and a version for Ubuntu
is available as a PPA and mirrored via i2p. To build the deb package, from the
root of this repository with the build dependencies installed(git, i2p, go,
debuild) run the command

        debuild -us -uc

to produce an unsigned deb for personal use only. For packagers,

        debuild -S

will produce a viable source package for use with Launchpad PPA's and other
similar systems.

### TODO

* Improve recovery on failed sockets
* Implement `STREAM FORWARD`
* Implement datagrams (Repliable and Anon)

