maketables: maketables.go maketesttables.go
	go build maketables.go
	go build maketesttables.go

tables:	maketables
	./maketables > tables.go
	gofmt -w tables.go
	./maketesttables > tables_test.go
	gofmt -w tables_test.go
