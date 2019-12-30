export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: shard plugin

clean:
	rm -rf ./libshared.so ./libshared.h ./libplugin.so
	@echo "Done clean"

shard:
	go build -buildmode=c-shared -o libshared.so  libshared.go
	@echo "Done build"

plugin:
	go build -buildmode=plugin -o libplugin.so  libplugin.go
	@echo "Done build"

.PHONY: clean
