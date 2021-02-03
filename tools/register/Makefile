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
PUB_KEY_FILE     := ../../keys/v1/pubkey.pem
PubKeyVAL :=$(shell cat $(PUB_KEY_FILE))

all: linux win64 mac

show:
	@echo $(PubKeyVAL)

clean:
	rm -rf ./register ./register.exe ./register.app ./license.dat
	@echo "Done clean"

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \
    " \
    -X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
    -X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
    -X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
    -X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
    -X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
    -X '${ECC_KEYS}.ECDSA_PUBLICKEY=$(PubKeyVAL)' \
    " \
	-o register main.go
	@echo "Done build"

win64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags \
    " \
    -X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
    -X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
    -X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
    -X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
    -X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
    -X '${ECC_KEYS}.ECDSA_PUBLICKEY=$(PubKeyVAL)' \
    " \
	-o register.exe main.go
	@echo "Done build"

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \
    " \
    -X '${VERSION_DIR}.buildVersion=${BUILD_VERSION}' \
    -X '${VERSION_DIR}.buildName=${BUILD_NAME}' \
    -X '${VERSION_DIR}.buildBranch=${BUILD_BRANCH}' \
    -X '${VERSION_DIR}.buildCommitID=${BUILD_COMMITID}' \
    -X '${VERSION_DIR}.buildTime=${BUILD_TIME}' \
    -X '${ECC_KEYS}.ECDSA_PUBLICKEY=$(PubKeyVAL)' \
    " \
	-o register.app main.go
	@echo "Done build"
 
.PHONY: clean
