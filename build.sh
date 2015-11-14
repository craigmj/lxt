#!/bin/bash
set -e
export GOPATH=`pwd`
for p in \
	"gopkg.in/lxc/go-lxc.v2" \
	"github.com/craigmj/commander" \
	"github.com/golang/glog" \
	; do
	if [ ! -d src/$p ]; then
		go get $p
	fi
done
if [ ! -d bin ]; then
	mkdir bin
fi
go build -o bin/lxt -a src/cmd/lxt.go
if [ ! `which lxt` ]; then 
	sudo ln -s `pwd`/bin/lxt /usr/bin/lxt
fi

