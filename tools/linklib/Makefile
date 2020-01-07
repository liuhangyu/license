export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: shard plugin

clean:
	rm -rf ./libshared.so ./libshared.h ./libplugin.so ./libshared.dylib
	@echo "Done clean"

shard:
ifeq ($(shell uname), Linux) 
	go build -buildmode=c-shared -o libshared.so  libshared.go
else ifeq ($(shell uname), Darwin) 
    go build -buildmode=c-shared -o libshared.dylib  libshared.go
endif
	@echo "Done build"

plugin:
	go build -buildmode=plugin -o libplugin.so  libplugin.go
	@echo "Done build"

.PHONY: clean
