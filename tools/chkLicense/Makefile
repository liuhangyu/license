export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: chkshard chkplugin

clean:
	rm -rf ./chkplugin ./chkshard ./license.log
	@echo "Done clean"

chkshard:
	gcc -o chkshard ./shard/main.c -ldl -lshared -L ../linklib -I ../linklib/
	@echo "Done build"

chkplugin:
	go build -o chkplugin ./plugin/main.go
	@echo "Done build"

run:
	./chkplugin -l ../cli  -lib ../linklib/libplugin.so -p switch-directory-chain
	./chkshard  "../cli"  "switch-directory-chain" "../linklib/libshared.so"

.PHONY: clean
