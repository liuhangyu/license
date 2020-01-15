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
	export LD_LIBRARY_PATH=../linklib
	./chkplugin -l ../register  -lib ../linklib/libplugin.so -p switch-directory-chain
	./chkshard  "../register"  "switch-directory-chain" "../linklib/libshared.so"

.PHONY: clean
