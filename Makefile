export GOPATH := $(PWD)/..

all: get
	go build .

get:
	go get -d .

release: get
	GOOS=darwin GOARCH=amd64 go build -o release/darwin_amd64/compose-repo .
	GOOS=linux GOARCH=amd64 go build -o release/linux_amd64/compose-repo .
	GOOS=linux GOARCH=386 go build -o release/linux_386/compose-repo .
