export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: chkshard chkplugin realcheck

clean:
	rm -rf ./chkplugin ./chkshard ./realcheck ./license.log
	@echo "Done clean"

chkshard:
	gcc -o chkshard ./shared/main.c -ldl -llicense -L ../../linklib/shared -I ../../linklib/shared
	@echo "Done build"

chkplugin:
	go build -o chkplugin ./plugin/main.go
	@echo "Done build"

realcheck:
	go build -o realcheck ./real/main.go

run:
	export LD_LIBRARY_PATH=../../linklib/shared
	./chkplugin -l ../register  -lib ../../linklib/plugin/liblicense.so -p switch-directory-chain
	./chkshard  "../register"  "switch-directory-chain" "../../linklib/shared/liblicense.so"
	./realcheck  -l ../register  -lib ../../linklib/plugin/liblicense.so -p switch-directory-chain

.PHONY: clean chkshard chkplugin realcheck
