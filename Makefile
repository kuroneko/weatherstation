# Makefile for wxstation
#
export GOPATH
GOPATH=$(CURDIR)

all:	bin/wxstation

ALLSOURCES=\
	src/probe/probe.go \
	src/wxstation/main.go \


bin/wxstation:	src/wxstation/main.go
	go install -a wxstation

deps:
	go get -d wxstation

clean:
	$(RM) bin/wxstation
	$(RM) -r pkg
