export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: linux win64

clean:
	rm -rf ./licensemgr ./licensemgr.exe
	@echo "Done clean"

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o licensemgr main.go
	@echo "Done build"

win64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o licensemgr.exe main.go
	@echo "Done build"

.PHONY: clean
