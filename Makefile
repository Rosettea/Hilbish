PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
LIBDIR ?= $(PREFIX)/share/hilbish

MY_GOFLAGS = -ldflags "-s -w"

all: dev

dev: MY_GOFLAGS = -ldflags "-s -w -X main.gitCommit=$(shell git rev-parse --short HEAD) -X main.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
dev: build

build:
	go build $(MY_GOFLAGS)

install:
	install -v -d "$(DESTDIR)$(BINDIR)/" && install -m 0755 -v hilbish "$(DESTDIR)$(BINDIR)/hilbish"
	mkdir -p "$(DESTDIR)$(LIBDIR)"
	cp -r libs docs emmyLuaDocs nature .hilbishrc.lua "$(DESTDIR)$(LIBDIR)"
	grep -qxF "$(DESTDIR)$(BINDIR)/hilbish" /etc/shells || echo "$(DESTDIR)$(BINDIR)/hilbish" >> /etc/shells

uninstall:
	rm -vrf \
			"$(DESTDIR)$(BINDIR)/hilbish" \
			"$(DESTDIR)$(LIBDIR)"
	sed -i '/hilbish/d' /etc/shells

clean:
	go clean

.PHONY: all dev build install uninstall clean