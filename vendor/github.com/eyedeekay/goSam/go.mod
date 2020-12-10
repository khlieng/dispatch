module github.com/eyedeekay/goSam

require (
	github.com/eyedeekay/sam3 v0.32.32
	github.com/getlantern/go-socks5 v0.0.0-20171114193258-79d4dd3e2db5
	github.com/getlantern/golog v0.0.0-20201105130739-9586b8bde3a9 // indirect
	github.com/getlantern/netx v0.0.0-20190110220209-9912de6f94fd // indirect
	github.com/getlantern/ops v0.0.0-20200403153110-8476b16edcd6 // indirect
)

//replace github.com/eyedeekay/gosam v0.1.1-0.20190814195658-27e786578944 => github.com/eyedeekay/goSam ./

replace github.com/eyedeekay/gosam v0.32.1 => ./

replace github.com/eyedeekay/goSam v0.32.1 => ./

go 1.13
