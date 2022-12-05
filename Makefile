default: test

test: 
	cd calculator
	go test -v ./...

build:
	cd cmd/calculator main
	go build -o calculator