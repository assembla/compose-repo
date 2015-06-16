export GOPATH := $(PWD)/..
all:
	go get -d .
	go build .
