default: test

test: 
	cd calculator && go test -v ./...

build:
	cd caculator && go build -o calculator
	