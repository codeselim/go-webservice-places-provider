.PHONY: help

help:
	echo "This is help task"

build:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o places-search places-search.go