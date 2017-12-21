SHELL := /bin/bash

NAME = gwatch
VERSION = 0.0.4
DESCTIPTION = "execute a program periodically, graphing the output fullscreen"
BUILD_ROOT = build
BUILD_REV = `git rev-parse HEAD`
BUILD_HOST = `hostname --long`
BUILD_DATE = `date`
PACKAGE_OUT = packages
BINARIES_OUT = binaries

# TODO use this like at all
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

default: clean test build package

test:
	go test ./...

# TODO: build a sane tree of output arch and goos binaries, do the same for
#       packaging. For now i'll support darwin and linux but I wouldn't want
#       to leave out the bsd people
build: test
	mkdir $(BINARIES_OUT)
	GOOS=linux go build -o $(BINARIES_OUT)/linux/gwatch
	GOOS=darwin go build -o $(BINARIES_OUT)/darwin/gwatch
	# GOOS=freebsd go build -o $(BINARIES_OUT)/freebsd/gwatch
	# GOOS=netbsd go build -o $(BINARIES_OUT)/netbsd/gwatch
	# GOOS=openbsd go build -o $(BINARIES_OUT)/openbsd/gwatch

install:
	# Debatably wise...
	go install

package: pre_package_linux package_deb package_rpm package_homebrew

pre_package_linux:
	mkdir -p $(BUILD_ROOT)/linux
	# Bin
	mkdir -p $(BUILD_ROOT)/linux/usr/bin
	cp $(BINARIES_OUT)/linux/gwatch $(BUILD_ROOT)/linux/usr/bin/
	chmod 0755 $(BUILD_ROOT)/linux/usr/bin/gwatch
	# Manual
	mkdir -p $(BUILD_ROOT)/linux/usr/share/man/man1
	cp gwatch.1 $(BUILD_ROOT)/linux/usr/share/man/man1/
	gzip $(BUILD_ROOT)/linux/usr/share/man/man1/gwatch.1
	# Package target
	mkdir $(PACKAGE_OUT)

package_deb:
	fpm \
	-s dir \
	-t deb \
	-C $(BUILD_ROOT)/linux \
	-p $(PACKAGE_OUT)/ \
	--name $(NAME) \
	--version $(VERSION) \
	--description $(DESCTIPTION) \
	--deb-no-default-config-files

package_rpm:
	fpm \
	-s dir \
	-t rpm \
	-C $(BUILD_ROOT)/linux \
	-p $(PACKAGE_OUT)/ \
	--name $(NAME) \
	--version $(VERSION) \
	--description $(DESCTIPTION)

package_homebrew:
	mkdir -p $(BUILD_ROOT)/darwin
	cp $(BINARIES_OUT)/darwin/gwatch $(BUILD_ROOT)/darwin/
	chmod 555 $(BUILD_ROOT)/darwin/gwatch
	cp gwatch.1 $(BUILD_ROOT)/darwin/
	gzip $(BUILD_ROOT)/darwin/gwatch.1
	tar -C $(BUILD_ROOT)/darwin -zcvf gwatch-$(VERSION).tar.gz gwatch gwatch.1.gz
	mv gwatch-$(VERSION).tar.gz packages


clean:
	rm -rf $(BUILD_ROOT)
	rm -rf $(PACKAGE_OUT)
	rm -rf $(BINARIES_OUT)
