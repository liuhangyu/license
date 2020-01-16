export WORKSPACE = $(shell pwd)
export WORKDIR=$(WORKSPACE)
export GO111MODULE=on


all: chkshard chkplugin gorealcheck crealcheck

clean:
	rm -rf ./chkplugin ./chkshard ./gorealcheck ./crealcheck ./license.log
	@echo "Done clean"

chkshard:
	gcc -o chkshard ./shared/main.c -ldl -llicense -L ../../linklib/shared -I ../../linklib/shared
	@echo "Done build"

chkplugin:
	go build -o chkplugin ./plugin/main.go
	@echo "Done build"

gorealcheck:
	go build -o gorealcheck ./goreal/main.go
	@echo "Done build"

crealcheck:
	gcc -o crealcheck ./creal/main.c -lpthread -ldl -llicense   -L ../../linklib/shared -I ../../linklib/shared
	@echo "Done build"

run:
	export LD_LIBRARY_PATH=../../linklib/shared
	./chkplugin -l ../register  -lib ../../linklib/plugin/liblicense.so -p switch-directory-chain
	./chkshard  "../register"  "switch-directory-chain" "../../linklib/shared/liblicense.so"
	./gorealcheck  -l ../register  -lib ../../linklib/plugin/liblicense.so -p switch-directory-chain

runc:
	export LD_LIBRARY_PATH=../../linklib/shared
	./crealcheck  "../register"  "switch-directory-chain" "../../linklib/shared/liblicense.so"


.PHONY: clean chkshard chkplugin gorealcheck crealcheck
