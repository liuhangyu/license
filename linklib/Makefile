export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on

#version
VERSION_DIR     := license/public
ECC_KEYS        := license/public

BUILD_VERSION   = $(shell git describe --abbrev=0 --tags)
BUILD_BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMITID  := $(shell git log --pretty=format:"%h" -1 )
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := license

#公私秘钥
PUB_KEY_FILE     := ../keys/v1/pubkey.pem
PubKeyVAL :=$(shell cat $(PUB_KEY_FILE))

all: shard plugin

show:
	@echo $(PubKeyVAL)

clean:
	rm -rf ./shared/liblicense.so ./shared/liblicense.h  ./shared/liblicense.dylib  ./plugin/liblicense.so
	@echo "Done clean"

shard:
ifeq ($(shell uname), Linux) 
	go build -ldflags \
    " \
    -X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
    -X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
    -X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
    -X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
    -X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
    -X '${ECC_KEYS}.ECDSA_PUBLICKEY=${PubKeyVAL}' \
    " \
	-buildmode=c-shared -o $(PWD)/shared/liblicense.so  $(PWD)/shared/libshared.go
else ifeq ($(shell uname), Darwin) 
	go build -ldflags \
    " \
    -X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
    -X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
    -X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
    -X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
    -X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
    -X '${ECC_KEYS}.ECDSA_PUBLICKEY=${PubKeyVAL}' \
    " \
	-buildmode=c-shared -o $(PWD)/shared/liblicense.dylib  $(PWD)/shared/libshared.go
endif
	@echo "Done build"

plugin:
	go build -ldflags \
	" \
	-X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
	-X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
	-X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
	-X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
	-X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
	-X '${ECC_KEYS}.ECDSA_PUBLICKEY=${PubKeyVAL}' \
	" \
    -buildmode=plugin -o  $(PWD)/plugin/liblicense.so $(PWD)/plugin/libplugin.go
	@echo "Done build"

.PHONY: clean all plugin shard show