SHELL := /bin/bash

NAME = gwatch
VERSION = 1.0.0
DESCTIPTION = "execute a program periodically, graphing the output fullscreen"
BUILD_ROOT = build
BUILD_REV = `git rev-parse HEAD`
BUILD_HOST = `hostname --long`
BUILD_DATE = `date`
PACKAGE_OUT = packages

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

all: test build install

test:
	go test ./...

build: test
	go build

install:
	go install

package: pre_package package_deb package_rpm

pre_package:
	mkdir $(BUILD_ROOT)
	# Bin
	mkdir -p $(BUILD_ROOT)/usr/bin
	cp gwatch $(BUILD_ROOT)/usr/bin/
	# Manual
	mkdir -p $(BUILD_ROOT)/usr/share/man/man1
	cp gwatch.1 $(BUILD_ROOT)/usr/share/man/man1/
	gzip $(BUILD_ROOT)/usr/share/man/man1/gwatch.1
	mkdir $(PACKAGE_OUT)

package_deb:
	fpm \
	-s dir \
	-t deb \
	-C $(BUILD_ROOT) \
	-p $(PACKAGE_OUT)/ \
	--name $(NAME) \
	--version $(VERSION) \
	--description $(DESCTIPTION) \
	--deb-no-default-config-files


package_rpm:
	fpm \
	-s dir \
	-t rpm \
	-C $(BUILD_ROOT) \
	-p $(PACKAGE_OUT)/ \
	--name $(NAME) \
	--version $(VERSION) \
	--description $(DESCTIPTION)

clean:
	rm -rf $(BUILD_ROOT)
	rm -rf $(PACKAGE_OUT)
	rm gwatch
