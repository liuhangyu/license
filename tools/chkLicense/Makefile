export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: chkshard chkplugin

clean:
	rm -rf ./chkplugin ./chkshard
	@echo "Done clean"

chkshard:
	gcc -o chkshard ./shard/main.c -ldl -lshared -L ../linklib -I ../linklib/
	@echo "Done build"

chkplugin:
	go build -o chkplugin ./plugin/main.go
	@echo "Done build"

run:
	./chkplugin -l ../cli/license.dat -lib ../linklib/libplugin.so -p switch-directory-chain
	./chkshard  "../cli/license.dat"  "switch-directory-chain" "../linklib/libshared.so"

.PHONY: clean
