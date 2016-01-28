configure:
	gb vendor update --all

build:
	gofmt -w src/seqrequest
	go tool vet src/seqrequest/*.go
	gb test
	gb build