CMD ?= sunc
PREFIX ?= /usr/local/bin

all: build

deps:
	go get github.com/syndtr/gocapability/capability

build: deps
	go build -o build/$(CMD)

install: build/$(CMD)
	cp build/$(CMD) $(PREFIX)/$(CMD)
	chmod u+s $(PREFIX)/$(CMD)

