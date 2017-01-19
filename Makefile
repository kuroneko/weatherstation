# Makefile for wxstation
#
export GOPATH
GOPATH=$(CURDIR)

all:	bin/wxstation

ALLSOURCES=\
	src/probe/probe.go \
	src/wxstation/main.go \
	src/wxstation/graph.go \
	src/wxstation/web.go \


bin/wxstation:	$(ALLSOURCES)
	go install -v wxstation

deps:
	go get -d wxstation

clean:
	$(RM) bin/wxstation
	$(RM) -r pkg
