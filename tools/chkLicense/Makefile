export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: chkshard chkplugin

clean:
	rm -rf ./chkplugin ./chkshard ./license.log
	@echo "Done clean"

chkshard:
	gcc -o chkshard ./shard/main.c -ldl -llicense -L ../linklib/shared -I ../linklib/shared
	@echo "Done build"

chkplugin:
	go build -o chkplugin ./plugin/main.go
	@echo "Done build"

run:
	export LD_LIBRARY_PATH=../linklib/shared
	./chkplugin -l ../register  -lib ../linklib/plugin/liblicense.so -p switch-directory-chain
	./chkshard  "../register"  "switch-directory-chain" "../linklib/shared/liblicense.so"

.PHONY: clean
