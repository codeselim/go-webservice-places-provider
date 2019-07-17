.PHONY: run

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

configure:	## Install & configure project
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

test:
	@go test ./...

build: fmt config test
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o places places.go

fmt:
	@go fmt ./...

run: fmt
	go run places.go
