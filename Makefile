export GOPATH := $(PWD)/..

all: get
	go build .

get:
	go get -d .

release: get
	GOOS=darwin GOARCH=amd64 go build -o release/darwin_amd64/compose-repo .
	cd release/darwin_amd64 && \
	zip -r compose-repo_darwin_amd64.zip compose-repo && \
	mv compose-repo_darwin_amd64.zip ..

	GOOS=linux GOARCH=amd64 go build -o release/linux_amd64/compose-repo .
	cd release/linux_amd64 && \
	zip -r compose-repo_linux_amd64.zip compose-repo && \
	mv compose-repo_linux_amd64.zip ..

	GOOS=linux GOARCH=386 go build -o release/linux_386/compose-repo .
	cd release/linux_386 && \
	zip -r compose-repo_linux_386.zip compose-repo && \
	mv compose-repo_linux_386.zip ..
